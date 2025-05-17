package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type LegalEntitySender struct {
	producer sarama.SyncProducer
	topic    string
}

type BankAccountSender struct {
	producer sarama.SyncProducer
	topic    string
}

type CreatedMessage struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewKafkaSyncProducer() sarama.SyncProducer {
	brokers := []string{"kafka:9092"}

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	return producer
}

func NewLegalEntitySender(producer sarama.SyncProducer) *LegalEntitySender {
	return &LegalEntitySender{producer: producer, topic: "legal-entities-created"}
}

func NewBankAccountSender(producer sarama.SyncProducer) *BankAccountSender {
	return &BankAccountSender{producer: producer, topic: "bank-accounts-created"}
}

func (s *LegalEntitySender) Send(ctx context.Context, id string, createdAt time.Time) error {
	return sendToKafka(ctx, s.producer, s.topic, id, createdAt)
}

func (s *BankAccountSender) Send(ctx context.Context, id string, createdAt time.Time) error {
	return sendToKafka(ctx, s.producer, s.topic, id, createdAt)
}

func sendToKafka(ctx context.Context, producer sarama.SyncProducer, topic, id string, createdAt time.Time) error {
	msg := CreatedMessage{
		ID:        id,
		CreatedAt: createdAt,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bytes),
	}
	_, _, err = producer.SendMessage(kafkaMsg)
	return err
}
