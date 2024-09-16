package producer

import (
	"go-kafka/internal/config"
	"go-kafka/internal/lib/logger/sl"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerService struct {
	log      *slog.Logger
	cfg      config.Kafka
	producer *kafka.Producer
}

func NewProducerService(cfg config.Kafka, log *slog.Logger) (*ProducerService, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaBrokers,
		"acks":              cfg.ProducerAsks,
	})
	if err != nil {
		log.Error("failed to create producer", sl.Err(err))
		return nil, err

	}
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error(
						"Delivery failed",
						slog.Any("topic", ev.TopicPartition.Topic),
						slog.Any("partition", ev.TopicPartition.Partition),
						sl.Err(ev.TopicPartition.Error),
					)
				} else {
					log.Info(
						"Message delivered successfully",
						slog.Any("topic", ev.TopicPartition.Topic),
						slog.Any("partition", ev.TopicPartition.Partition),
						slog.String("message", string(ev.Value)),
					)
				}
			}
		}
	}()

	return &ProducerService{
		log:      log,
		cfg:      cfg,
		producer: producer,
	}, nil
}

func (p *ProducerService) StopProducer() {
	p.log.Info("Stopping producer")
	p.producer.Close()
}

func (p *ProducerService) ProduceMessage(value string, topic string) error {
	const op = "services.producer.ProduceMessage"

	p.log = slog.With(
		slog.String("op", op),
	)

	var err error
	for range p.cfg.ProduceRetries {
		err = p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}, nil)
		if err != nil {
			p.log.Error("failed to produce message", sl.Err(err))
			continue
		}
		break
	}

	if err != nil {
		p.log.Error("failed to produce message after all retries", sl.Err(err))
		return err
	}

	p.log.Info("message successfully sent")
	return nil
}
