package kafka

import (
	"os"

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

func (in *Instance) NewReader(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{in.endpoint},
		Topic:   topic,
		// GroupID:   groupID,
		Partition: 0,
		MaxBytes:  10e6,
	})
}
