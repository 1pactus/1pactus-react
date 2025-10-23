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
	KafkaTopicBlocks = "onepacd-blocks"
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
	s.blocksReader = store.GetReader(KafkaTopicBlocks)
	s.writer = store.GetWriter()
}

func (s *kafkaStore) Topics() []string {
	return []string{KafkaTopicBlocks}
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
		Topic: KafkaTopicBlocks,
		Key:   nil,
		Value: data,
		Time:  time.Now(),
	}

	return s.writer.WriteMessages(ctx, message)
}

func (s *kafkaStore) ConsumeMessages(topic string, handler func(kafka.Message) error) error {
	reader := s.Kafka.GetReader(topic)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), s.Kafka.GetTimeout())

		message, err := reader.ReadMessage(ctx)
		cancel()

		if err != nil {
			return fmt.Errorf("failed to read message: %v", err)
		}

		if err := handler(message); err != nil {
			log.Printf("Error handling message: %v", err)
			continue
		}
	}
}

func (s *kafkaStore) GetLastBlockHeight() (int64, error) {
	partitionsLastMessage, err := s.Kafka.GetAllPartitionsLastMessage(KafkaTopicBlocks)

	if err != nil {
		return 0, fmt.Errorf("GetAllPartitionsLastMessage failed: %w", err)
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

	if height == -1 {
		return 0, fmt.Errorf("no blocks found in topic %s", KafkaTopicBlocks)
	}

	return height, nil
}
