package main

import (
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/bsergik/tough-dev/services/auth/internal/database/psql"
	"github.com/bsergik/tough-dev/services/auth/internal/message-broker/kafka"
	"github.com/bsergik/tough-dev/services/auth/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	keycloakAddress, ok := os.LookupEnv("KEYCLOAK_ADDRESS")
	if !ok {
		log.Panic().Msg("Env KEYCLOAK_ADDRESS is required")
	}

	db := psql.NewPsql()
	mq := kafka.NewMessageBroker()
	cloak := gocloak.NewClient(keycloakAddress)
	srv := server.NewServer(cloak, db, mq)

	bindAddr, ok := os.LookupEnv("SERVER_BIND_ADDRESS")
	if !ok {
		log.Panic().Msgf("Env SERVER_BIND_ADDRESS is required")
	}
	srv.Start(bindAddr)
}
