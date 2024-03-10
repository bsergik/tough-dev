package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bsergik/tough-dev/services/inventory/internal/message-broker/model"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Instance struct {
	endpoint string
}

func NewMessageBroker() *Instance {
	endpoint, ok := os.LookupEnv("MQ_ENDPOINT")
	if !ok {
		log.Fatal().Msg("MQ_ENDPOINT is required")
	}

	return &Instance{
		endpoint: endpoint,
	}
}

func (in *Instance) newWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(in.endpoint),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: false,
		Logger:                 &log.Logger,
		ErrorLogger:            &log.Logger,
		// Transport: &kafka.Transport{
		// 	TLS: &tls.Config{
		// 		MinVersion:         tls.VersionTLS12,
		// 		InsecureSkipVerify: true,
		// 	},
		// },
	}
}

func (in *Instance) newReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{in.endpoint},
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
	})
}

func (in *Instance) PublishTaskCreated(ctx context.Context, task *model.Task) error {
	w := in.newWriter("task-created")
	defer w.Close()

	b, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	err = w.WriteMessages(ctx, kafka.Message{
		Value: b,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
