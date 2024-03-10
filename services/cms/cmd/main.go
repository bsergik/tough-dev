package main

import (
	"os"

	"github.com/bsergik/tough-dev/services/cms/internal/database/psql"
	"github.com/bsergik/tough-dev/services/cms/internal/message-broker/kafka"
	"github.com/bsergik/tough-dev/services/cms/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	db := psql.NewPsql()
	mq := kafka.NewMessageBroker()
	srv := server.NewServer(db, mq)

	bindAddr, ok := os.LookupEnv("SERVER_BIND_ADDRESS")
	if !ok {
		log.Panic().Msgf("Env SERVER_BIND_ADDRESS is required")
	}
	srv.Start(bindAddr)
}
