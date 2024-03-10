package service

import (
	"context"
	"time"

	"github.com/bsergik/tough-dev/services/cms/internal/database/psql"
	"github.com/bsergik/tough-dev/services/cms/internal/message-broker/kafka"
	kafkamodelv1 "github.com/bsergik/tough-dev/services/cms/pkg/message-broker/kafka/model/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Instance struct {
	mq *kafka.Instance
	db *psql.Instance
}

func NewService(mq *kafka.Instance, db *psql.Instance) *Instance {
	return &Instance{
		mq: mq,
		db: db,
	}
}

func (in *Instance) Start(ctx context.Context) {
	mqReader := in.mq.NewReader("user-stream", "user-group")
	defer mqReader.Close()

	offset, err := in.db.GetMQOffset(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get offset")
		return
	}

	mqReader.SetOffset(offset.Offset)

	for {
		m, err := mqReader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to read message from mq")
			continue
		}

		userV1 := &kafkamodelv1.User{}

		err = proto.Unmarshal(m.Value, userV1)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal user")
			continue
		}

		_, err = in.db.CreateUser(
			ctx,
			int(userV1.RoleId),
			userV1.PublicId,
			userV1.Email,
			userV1.Username,
			userV1.FirstName,
			userV1.LastName,
			userV1.Enabled,
			time.Now(),
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to save user")
			continue
		}
	}
}

func (in *Instance) saveFailedMessage(ctx context.Context, msg []byte) {
	if err := in.db.SaveFailedMessage(ctx, msg); err != nil {
		log.Error().Err(err).Msg("failed to save failed message")
	}
}
