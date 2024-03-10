package psql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bsergik/tough-dev/services/cms/internal/database/model"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type Instance struct {
	conn *pgx.Conn
}

func NewPsql() *Instance {
	return &Instance{}
}

func (in *Instance) Connect(ctx context.Context) error {
	host, ok := os.LookupEnv("PSQL_HOST")
	if !ok {
		log.Error().Msg("PSQL_HOST is required")
		return fmt.Errorf("PSQL_HOST is required")
	}

	port, ok := os.LookupEnv("PSQL_PORT")
	if !ok {
		log.Error().Msg("PSQL_PORT is required")
		return fmt.Errorf("PSQL_PORT is required")
	}

	user, ok := os.LookupEnv("PSQL_USER")
	if !ok {
		log.Error().Msg("PSQL_USER is required")
		return fmt.Errorf("PSQL_USER is required")
	}

	pass, ok := os.LookupEnv("PSQL_PASS")
	if !ok {
		log.Error().Msg("PSQL_PASS is required")
		return fmt.Errorf("PSQL_PASS is required")
	}

	dbname, ok := os.LookupEnv("PSQL_DBNAME")
	if !ok {
		log.Error().Msg("PSQL_DBNAME is required")
		return fmt.Errorf("PSQL_DBNAME is required")
	}

	psqlURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)

	conn, err := pgx.Connect(ctx, psqlURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to psql")
		return fmt.Errorf("failed to connect to psql: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		conn.Close(ctx)
		log.Error().Err(err).Msg("failed to ping psql")
		return fmt.Errorf("failed to ping psql: %w", err)
	}

	go func() {
		<-ctx.Done()
		conn.Close(ctx)
	}()

	in.conn = conn

	return nil
}

func (in *Instance) GetMQOffset(ctx context.Context) (*model.MQOffset, error) {
	var offset model.MQOffset
	err := in.conn.QueryRow(ctx, "SELECT offset, updated_at FROM mq_offset ORDER BY DESC updated_at LIMIT 1").Scan(&offset.Offset)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &offset, nil
		}

		log.Error().Err(err).Msg("failed to get mq offset")
		return nil, fmt.Errorf("failed to get mq offset: %w", err)
	}

	return &offset, nil
}

func (in *Instance) SaveFailedMessage(ctx context.Context, msg []byte) error {
	_, err := in.conn.Exec(ctx, "INSERT INTO failed_messages (topic, group_id, message, received_at) VALUES ($1, $2, $3, $4)", msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to save failed message")
		return fmt.Errorf("failed to save failed message: %w", err)
	}

	return nil
}

func (in *Instance) CreateUser(
	ctx context.Context,
	roleID int,
	publicID string,
	email string,
	username string,
	firstName string,
	lastName string,
	enabled bool,
	createdAt time.Time,
) (*model.User, error) {
	_, err := in.conn.Exec(ctx, "INSERT INTO users "+
		"(role_id, public_id, email, username, first_name, last_name, enabled, createdAt) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		roleID, publicID, email, username, firstName, lastName, enabled, createdAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user := model.User{}
	err = in.conn.QueryRow(ctx, "SELECT id, role_id, public_id, email, username, first_name, last_name, enabled, created_at, updated_at "+
		"FROM users WHERE public_id = $1", publicID).
		Scan(&user.ID, &user.RoleID, &user.PublicID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Enabled, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
