package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/1pactus/1pactus-react/config"
	"github.com/1pactus/1pactus-react/store/storedriver"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const (
	kafkaTopicBlocks = "onepacd-blocks"
)

type kafkaStore struct {
	storedriver.Kafka
	conf         *config.KafkaConfig
	blocksReader *kafka.Reader
	writer       *kafka.Writer
}

func (s *kafkaStore) Init(store storedriver.Kafka, conf *config.KafkaConfig) {
	s.Kafka = store
	s.conf = conf
	s.blocksReader = store.GetReader(kafkaTopicBlocks)
	s.writer = store.GetWriter()
}

func (s *kafkaStore) Topics() []string {
	return []string{kafkaTopicBlocks}
}

func (s *kafkaStore) SendMessage(topic, key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.Kafka.GetTimeout())
	defer cancel()

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	return s.writer.WriteMessages(ctx, message)
}

func (s *kafkaStore) SendBlock(block *pactus.GetBlockResponse) error {
	data, err := proto.Marshal(block)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.Kafka.GetTimeout())
	defer cancel()

	message := kafka.Message{
		Topic: kafkaTopicBlocks,
		Key:   nil,
		Value: data,
		Time:  time.Now(),
	}

	return s.writer.WriteMessages(ctx, message)
}

func (s *kafkaStore) ConsumeBlocks(groupID string, offset int64, handler func(*pactus.GetBlockResponse) (bool, error)) error {
	reader := s.Kafka.GetReader(kafkaTopicBlocks, storedriver.NewReaderOptions().
		WithGroupID(groupID).
		WithSeekOffset(offset))

	defer reader.Close()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), s.Kafka.GetTimeout())

		message, err := reader.ReadMessage(ctx)
		cancel()

		if err != nil {
			return fmt.Errorf("failed to read message: %v", err)
		}

		var block pactus.GetBlockResponse
		if err := proto.Unmarshal(message.Value, &block); err != nil {
			return err
		}

		if ok, err := handler(&block); err != nil {
			return err
		} else {
			if !ok {
				// manual stop
				break
			}
		}
	}

	return nil
}

var ErrorKafkaTopicEmpty = fmt.Errorf("kafka topic is empty")

func (s *kafkaStore) GetLastBlockHeight() (int64, error) {
	partitionsLastMessage, err := s.Kafka.GetAllPartitionsLastMessage(kafkaTopicBlocks)

	if err != nil {
		return 0, fmt.Errorf("GetAllPartitionsLastMessage failed: %w", err)
	}

	if len(partitionsLastMessage) == 0 {
		return 0, ErrorKafkaTopicEmpty
	}

	height := int64(-1)

	for _, message := range partitionsLastMessage {
		var block pactus.GetBlockResponse
		if err := proto.Unmarshal(message.Value, &block); err != nil {
			return 0, fmt.Errorf("failed to unmarshal block: %w", err)
		}

		if int64(block.Height) > height {
			height = int64(block.Height)
		}
	}

	return height, nil
}

func (s *kafkaStore) GetBlockHeightOffset(height int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Kafka.GetTimeout())
	defer cancel()

	offset, err := s.Kafka.FindOffset(ctx, kafkaTopicBlocks, func(message kafka.Message) (int, error) {
		var block pactus.GetBlockResponse
		if err := proto.Unmarshal(message.Value, &block); err != nil {
			return 0, fmt.Errorf("failed to unmarshal block: %w", err)
		}

		if int64(block.Height) < height {
			return -1, nil
		} else if int64(block.Height) > height {
			return 1, nil
		}

		return 0, nil
	})

	if err != nil {
		return 0, fmt.Errorf("FindOffset failed: %w", err)
	}

	return offset, nil
}
