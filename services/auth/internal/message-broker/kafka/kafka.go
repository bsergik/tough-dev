package kafka

import (
	"context"
	"fmt"
	"os"

	kafkamodelv1 "github.com/bsergik/tough-dev/services/auth/pkg/message-broker/kafka/model/v1"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
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

func (in *Instance) PublishUserCreated(ctx context.Context, user *kafkamodelv1.User) error {
	w := in.newWriter("user-stream")
	defer w.Close()

	b, err := proto.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	err = w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("created"),
		Value: b,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
