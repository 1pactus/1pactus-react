package storedriver

import (
	"context"
	"fmt"
	"time"

	"github.com/1pactus/1pactus-react/config"
	"github.com/1pactus/1pactus-react/log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

// ReaderOptions holds configuration options for creating a Kafka reader
type ReaderOptions struct {
	groupID     string
	startOffset int64
	seekOffset  *int64 // For specific numeric offsets that require Seek()
	minBytes    int
	maxBytes    int
	partition   *int
}

// NewReaderOptions creates a new ReaderOptions with default values
func NewReaderOptions() *ReaderOptions {
	return &ReaderOptions{
		startOffset: kafka.LastOffset,
		minBytes:    10e3, // 10KB
		maxBytes:    10e6, // 10MB
	}
}

// WithGroupID sets the consumer group ID
func (r *ReaderOptions) WithGroupID(groupID string) *ReaderOptions {
	r.groupID = groupID
	return r
}

// WithSeekOffset sets a specific numeric offset that will be applied using Seek()
func (r *ReaderOptions) WithSeekOffset(offset int64) *ReaderOptions {
	r.seekOffset = &offset
	r.startOffset = kafka.FirstOffset // Start from beginning to allow seeking
	return r
}

// WithFirstOffset sets the start offset to the beginning of the topic
func (r *ReaderOptions) WithFirstOffset() *ReaderOptions {
	r.startOffset = kafka.FirstOffset
	return r
}

// WithLastOffset sets the start offset to the end of the topic (latest messages)
func (r *ReaderOptions) WithLastOffset() *ReaderOptions {
	r.startOffset = kafka.LastOffset
	return r
}

// WithMinBytes sets the minimum number of bytes to fetch in each request
func (r *ReaderOptions) WithMinBytes(minBytes int) *ReaderOptions {
	r.minBytes = minBytes
	return r
}

// WithMaxBytes sets the maximum number of bytes to fetch in each request
func (r *ReaderOptions) WithMaxBytes(maxBytes int) *ReaderOptions {
	r.maxBytes = maxBytes
	return r
}

// WithPartition sets a specific partition to read from
func (r *ReaderOptions) WithPartition(partition int) *ReaderOptions {
	r.partition = &partition
	return r
}

// ReaderConfig holds configuration options for creating a Kafka reader (legacy)
type ReaderConfig struct {
	Topic       string
	GroupID     string
	StartOffset int64
	MinBytes    int
	MaxBytes    int
	Partition   *int // Optional: specific partition to read from
}

type IKafkaStore interface {
	Init(store Kafka, conf *config.KafkaConfig)
	Topics() []string
}

type Kafka interface {
	GetReader(topic string, options ...*ReaderOptions) *kafka.Reader
	GetWriter() *kafka.Writer
	GetConn() *kafka.Conn
	GetTimeout() time.Duration
	GetAllPartitionsLastMessage(topic string) (map[int]*kafka.Message, error)
	FindOffset(ctx context.Context, topic string, handler func(kafka.Message) (int, error)) (int64, error)
	GetLogger() log.ILogger
}

type kafkaImpl struct {
	conf    *config.KafkaConfig
	conn    *kafka.Conn
	writer  *kafka.Writer
	stores  []IKafkaStore
	timeout time.Duration
	log     log.ILogger
}

func (k *kafkaImpl) Close() {
	if k.conn != nil {
		k.conn.Close()
	}

	if k.writer != nil {
		k.writer.Close()
	}
}

func KafkaStart(name string, conf *config.KafkaConfig, stores []IKafkaStore) error {
	k := &kafkaImpl{
		conf:    conf,
		stores:  stores,
		timeout: time.Second * 10, // Default timeout
		log:     log.WithKv("module", "store").WithKv("kafka", name),
	}

	var err error

	maxRetry := 10

	for {
		if err = k.connect(); err == nil {
			k.initStores()

			k.log.Infof("kafka connect and initialized success")

			if err := k.ensureTopics(); err != nil {
				return fmt.Errorf("kafka [%s] ensure topics failed: %v", name, err)
			}
			go k.monitorConnection()
			return nil
		} else {
			maxRetry -= 1

			if maxRetry <= 0 {
				return fmt.Errorf("kafka [%s] connect failed after retries: %v", name, err)
			}

			k.log.Errorf("kafka [%s] connect failed: %v", name, err)

			time.Sleep(time.Second * 5)
		}
	}
}

func (k *kafkaImpl) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse brokers
	if len(k.conf.Brokers) == 0 {
		return fmt.Errorf("no kafka brokers configured")
	}

	// Create connection to the first broker for testing connectivity
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Configure SASL authentication if provided
	if k.conf.SASL.Enabled {
		switch k.conf.SASL.Mechanism {
		case "PLAIN":
			dialer.SASLMechanism = plain.Mechanism{
				Username: k.conf.SASL.Username,
				Password: k.conf.SASL.Password,
			}
		case "SCRAM-SHA-256":
			mechanism, err := scram.Mechanism(scram.SHA256, k.conf.SASL.Username, k.conf.SASL.Password)
			if err != nil {
				return fmt.Errorf("failed to create SCRAM mechanism: %v", err)
			}
			dialer.SASLMechanism = mechanism
		case "SCRAM-SHA-512":
			mechanism, err := scram.Mechanism(scram.SHA512, k.conf.SASL.Username, k.conf.SASL.Password)
			if err != nil {
				return fmt.Errorf("failed to create SCRAM mechanism: %v", err)
			}
			dialer.SASLMechanism = mechanism
		}
	}

	// Test connection
	conn, err := dialer.DialContext(ctx, "tcp", k.conf.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to kafka broker: %v", err)
	}

	k.conn = conn
	k.setupWriter(dialer)

	return nil
}

func (k *kafkaImpl) setupWriter(dialer *kafka.Dialer) {
	k.writer = &kafka.Writer{
		Addr:                   kafka.TCP(k.conf.Brokers...),
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           k.timeout,
		ReadTimeout:            k.timeout,
		RequiredAcks:           kafka.RequireOne,
		Async:                  true,
		AllowAutoTopicCreation: k.conf.AutoCreateTopics,
		Transport:              &kafka.Transport{Dial: dialer.DialFunc},
		BatchSize:              200,
		BatchTimeout:           time.Second * 2,
		Compression:            kafka.Snappy,
	}
}

func (k *kafkaImpl) initStores() {
	for _, store := range k.stores {
		store.Init(k, k.conf)
	}
}

func (k *kafkaImpl) ensureTopics() error {
	allTopics := make(map[string]bool)

	// Collect all topics from stores
	for _, store := range k.stores {
		topics := store.Topics()
		for _, topic := range topics {
			allTopics[topic] = true
		}
	}

	if len(allTopics) == 0 {
		return nil
	}

	// Get existing topics
	partitions, err := k.conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("failed to read partitions: %v", err)
	}

	existingTopics := make(map[string]bool)
	for _, partition := range partitions {
		existingTopics[partition.Topic] = true
	}

	// Get available brokers to determine max replication factor
	brokers, err := k.conn.Brokers()
	if err != nil {
		return fmt.Errorf("failed to get brokers: %v", err)
	}

	maxReplicationFactor := len(brokers)
	replicationFactor := k.conf.DefaultReplicationFactor

	// Ensure replication factor doesn't exceed available brokers
	if replicationFactor > maxReplicationFactor {
		replicationFactor = maxReplicationFactor
		k.log.Warnf("replication factor %d exceeds available brokers %d, using %d",
			k.conf.DefaultReplicationFactor, maxReplicationFactor, replicationFactor)
	}

	// Create missing topics
	var topicsToCreate []kafka.TopicConfig
	for topic := range allTopics {
		if !existingTopics[topic] {
			topicsToCreate = append(topicsToCreate, kafka.TopicConfig{
				Topic:             topic,
				NumPartitions:     k.conf.DefaultPartitions,
				ReplicationFactor: replicationFactor,
			})
		}
	}

	if len(topicsToCreate) > 0 {
		if err := k.conn.CreateTopics(topicsToCreate...); err != nil {
			return fmt.Errorf("failed to create topics: %v", err)
		}
		k.log.Infof("created %d topics with replication factor %d", len(topicsToCreate), replicationFactor)
	}

	return nil
}

func (k *kafkaImpl) monitorConnection() {
	healthcheck := time.Duration(k.conf.Healthcheck)

	for {
		func() {
			// Test connection by reading brokers
			_, err := k.conn.Brokers()
			if err != nil {
				k.log.Errorf("lost kafka connection, retrying: %v", err)

				for {
					if err := k.connect(); err == nil {
						k.initStores()
						k.log.Info("successfully reconnected to kafka")
						break
					}
					k.log.Errorf("error connecting to kafka: %v", err)
					time.Sleep(5 * time.Second)
				}
			}
		}()

		time.Sleep(healthcheck * time.Second)
	}
}

func (k *kafkaImpl) GetLogger() log.ILogger {
	return k.log
}

func (k *kafkaImpl) GetReader(topic string, options ...*ReaderOptions) *kafka.Reader {
	// Use default options if none provided
	var opts *ReaderOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = NewReaderOptions()
	}

	// Set default groupID if not specified
	groupID := opts.groupID
	if groupID == "" {
		groupID = k.conf.ConsumerGroup
	}

	if opts.seekOffset != nil {
		groupID = ""
	}

	// Create new reader
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Configure SASL for reader if enabled
	if k.conf.SASL.Enabled {
		switch k.conf.SASL.Mechanism {
		case "PLAIN":
			dialer.SASLMechanism = plain.Mechanism{
				Username: k.conf.SASL.Username,
				Password: k.conf.SASL.Password,
			}
		case "SCRAM-SHA-256":
			mechanism, err := scram.Mechanism(scram.SHA256, k.conf.SASL.Username, k.conf.SASL.Password)
			if err != nil {
				k.log.Errorf("failed to create SCRAM mechanism for reader: %v", err)
				return nil
			}
			dialer.SASLMechanism = mechanism
		case "SCRAM-SHA-512":
			mechanism, err := scram.Mechanism(scram.SHA512, k.conf.SASL.Username, k.conf.SASL.Password)
			if err != nil {
				k.log.Errorf("failed to create SCRAM mechanism for reader: %v", err)
				return nil
			}
			dialer.SASLMechanism = mechanism
		}
	}

	readerConfig := kafka.ReaderConfig{
		Brokers:          k.conf.Brokers,
		Topic:            topic,
		GroupID:          groupID,
		MinBytes:         opts.minBytes,
		MaxBytes:         opts.maxBytes,
		CommitInterval:   time.Second,
		Dialer:           dialer,
		StartOffset:      opts.startOffset,
		ReadBatchTimeout: time.Second * 5,
	}

	// Set specific partition if provided
	if opts.partition != nil {
		readerConfig.Partition = *opts.partition
		readerConfig.GroupID = "" // When reading from specific partition, GroupID should be empty
	}

	reader := kafka.NewReader(readerConfig)

	// If seekOffset is specified, seek to that position
	if opts.seekOffset != nil {
		if err := reader.SetOffset(*opts.seekOffset); err != nil {
			k.log.Errorf("failed to seek to offset %d: %v", *opts.seekOffset, err)
			reader.Close()
			return nil
		}
	}

	return reader
}

func (k *kafkaImpl) GetWriter() *kafka.Writer {
	return k.writer
}

func (k *kafkaImpl) GetConn() *kafka.Conn {
	return k.conn
}

func (k *kafkaImpl) GetTimeout() time.Duration {
	return k.timeout
}

func (k *kafkaImpl) GetAllPartitionsLastMessage(topic string) (map[int]*kafka.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), k.timeout)
	defer cancel()

	// Get all partitions for the topic
	partitions, err := k.conn.ReadPartitions(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions for topic %s: %v", topic, err)
	}

	result := make(map[int]*kafka.Message)

	// For each partition, get the last message
	for _, partition := range partitions {
		partitionID := partition.ID

		// Create a reader for this specific partition
		dialer := &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		}

		// Configure SASL for reader if enabled
		if k.conf.SASL.Enabled {
			switch k.conf.SASL.Mechanism {
			case "PLAIN":
				dialer.SASLMechanism = plain.Mechanism{
					Username: k.conf.SASL.Username,
					Password: k.conf.SASL.Password,
				}
			case "SCRAM-SHA-256":
				mechanism, err := scram.Mechanism(scram.SHA256, k.conf.SASL.Username, k.conf.SASL.Password)
				if err != nil {
					k.log.Errorf("failed to create SCRAM mechanism for partition reader: %v", err)
					continue
				}
				dialer.SASLMechanism = mechanism
			case "SCRAM-SHA-512":
				mechanism, err := scram.Mechanism(scram.SHA512, k.conf.SASL.Username, k.conf.SASL.Password)
				if err != nil {
					k.log.Errorf("failed to create SCRAM mechanism for partition reader: %v", err)
					continue
				}
				dialer.SASLMechanism = mechanism
			}
		}

		// Create a connection for this partition
		partitionConn, err := dialer.DialLeader(ctx, "tcp", k.conf.Brokers[0], topic, partitionID)
		if err != nil {
			k.log.Errorf("failed to dial leader for partition %d: %v", partitionID, err)
			continue
		}

		// Get the latest offset (high water mark)
		_, high, err := partitionConn.ReadOffsets()
		if err != nil {
			partitionConn.Close()
			k.log.Errorf("failed to read offsets for partition %d: %v", partitionID, err)
			continue
		}

		// If partition is empty, skip it
		if high <= 0 {
			partitionConn.Close()
			continue
		}

		// Set read position to the last message (high - 1)
		lastOffset := high - 1
		_, err = partitionConn.Seek(lastOffset, kafka.SeekAbsolute)
		if err != nil {
			partitionConn.Close()
			k.log.Errorf("failed to seek to last message in partition %d: %v", partitionID, err)
			continue
		}

		// Read the last message
		partitionConn.SetReadDeadline(time.Now().Add(k.timeout))
		message, err := partitionConn.ReadMessage(10e6) // 10MB max message size
		if err != nil {
			partitionConn.Close()
			k.log.Errorf("failed to read last message from partition %d: %v", partitionID, err)
			continue
		}

		// Convert to kafka.Message format
		kafkaMsg := &kafka.Message{
			Topic:     message.Topic,
			Partition: message.Partition,
			Offset:    message.Offset,
			Key:       message.Key,
			Value:     message.Value,
			Headers:   make([]kafka.Header, len(message.Headers)),
			Time:      message.Time,
		}

		// Convert headers
		for i, header := range message.Headers {
			kafkaMsg.Headers[i] = kafka.Header{
				Key:   header.Key,
				Value: header.Value,
			}
		}

		result[partitionID] = kafkaMsg
		partitionConn.Close()
	}

	return result, nil
}

func (k *kafkaImpl) FindOffset(ctx context.Context, topic string, handler func(kafka.Message) (int, error)) (int64, error) {
	// Get partitions for the topic
	partitions, err := k.conn.ReadPartitions(topic)
	if err != nil {
		return -1, fmt.Errorf("failed to read partitions for topic %s: %v", topic, err)
	}

	// Only handle single partition case
	if len(partitions) != 1 {
		return -1, fmt.Errorf("topic %s has %d partitions, only single partition is supported", topic, len(partitions))
	}

	partition := partitions[0]
	partitionID := partition.ID

	// Create dialer with SASL configuration
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Configure SASL for reader if enabled
	if k.conf.SASL.Enabled {
		switch k.conf.SASL.Mechanism {
		case "PLAIN":
			dialer.SASLMechanism = plain.Mechanism{
				Username: k.conf.SASL.Username,
				Password: k.conf.SASL.Password,
			}
		case "SCRAM-SHA-256":
			mechanism, err := scram.Mechanism(scram.SHA256, k.conf.SASL.Username, k.conf.SASL.Password)
			if err != nil {
				return -1, fmt.Errorf("failed to create SCRAM mechanism: %v", err)
			}
			dialer.SASLMechanism = mechanism
		case "SCRAM-SHA-512":
			mechanism, err := scram.Mechanism(scram.SHA512, k.conf.SASL.Username, k.conf.SASL.Password)
			if err != nil {
				return -1, fmt.Errorf("failed to create SCRAM mechanism: %v", err)
			}
			dialer.SASLMechanism = mechanism
		}
	}

	// Create connection for the partition
	partitionConn, err := dialer.DialLeader(ctx, "tcp", k.conf.Brokers[0], topic, partitionID)
	if err != nil {
		return -1, fmt.Errorf("failed to dial leader for partition %d: %v", partitionID, err)
	}
	defer partitionConn.Close()

	// Get offset range
	low, high, err := partitionConn.ReadOffsets()
	if err != nil {
		return -1, fmt.Errorf("failed to read offsets: %v", err)
	}

	// If partition is empty
	if high <= low {
		return -1, fmt.Errorf("partition is empty")
	}

	// Binary search implementation
	left := low
	right := high - 1
	resultOffset := int64(-1)

	for left <= right {
		mid := left + (right-left)/2

		// Seek to mid offset
		_, err = partitionConn.Seek(mid, kafka.SeekAbsolute)
		if err != nil {
			return -1, fmt.Errorf("failed to seek to offset %d: %v", mid, err)
		}

		// Read message at mid offset
		partitionConn.SetReadDeadline(time.Now().Add(k.timeout))
		message, err := partitionConn.ReadMessage(10e6) // 10MB max message size
		if err != nil {
			return -1, fmt.Errorf("failed to read message at offset %d: %v", mid, err)
		}

		// Convert to kafka.Message format
		kafkaMsg := kafka.Message{
			Topic:     message.Topic,
			Partition: message.Partition,
			Offset:    message.Offset,
			Key:       message.Key,
			Value:     message.Value,
			Headers:   make([]kafka.Header, len(message.Headers)),
			Time:      message.Time,
		}

		// Convert headers
		for i, header := range message.Headers {
			kafkaMsg.Headers[i] = kafka.Header{
				Key:   header.Key,
				Value: header.Value,
			}
		}

		// Call handler to compare the current message with target
		// Return values: < 0 if current < target, 0 if equal, > 0 if current > target
		cmp, err := handler(kafkaMsg)
		if err != nil {
			return -1, fmt.Errorf("handler error at offset %d: %v", mid, err)
		}

		if cmp == 0 {
			// Found exact match, but continue searching for the first occurrence (leftmost)
			resultOffset = mid
			right = mid - 1
		} else if cmp < 0 {
			// Current message is less than target, search right half
			left = mid + 1
		} else {
			// Current message is greater than target, search left half
			right = mid - 1
		}
	}

	if resultOffset == -1 {
		return -1, fmt.Errorf("target message not found in topic %s", topic)
	}

	return resultOffset, nil
}
