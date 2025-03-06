package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	"go.uber.org/zap"
)

type Producer interface {
	Produce(ctx context.Context, topic string, message []byte) error
}

type producer struct {
	saramaSyncProducer sarama.SyncProducer
	l                  logger.Interface
}

func NewProducer(cfg config.Config, l logger.Interface) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 1
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ClientID = cfg.KafkaClientId
	config.Metadata.Full = true

	saramaSyncProducer, err := sarama.NewSyncProducer([]string{cfg.KafkaBrokers}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}

	return &producer{
		saramaSyncProducer: saramaSyncProducer,
		l:                  l,
	}, nil
}

func (p producer) Produce(ctx context.Context, topic string, message []byte) error {
	p.l.Info("Kafka produce event", zap.String("Topic: ", topic), zap.ByteString("Message: ", message))

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Send message
	partition, offset, err := p.saramaSyncProducer.SendMessage(&msg)
	if err != nil {
		p.l.Error("Kafka produce event", zap.String("Error: ", err.Error()))
		return err
	}

	p.l.Info("Kafka produce event successfully", zap.Int32("Partition: ", partition), zap.Int64("Offset: ", offset))

	return nil
}
