package psql

import (
	"context"
	"fmt"
	"os"

	"github.com/bsergik/tough-dev/services/inventory/internal/database/model"
	"github.com/gofrs/uuid"
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

func (in *Instance) CreateTask(ctx context.Context, publicID, userID uuid.UUID, title, description string) (*model.Task, model.TaskStatus, error) {
	tx, err := in.conn.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
		return nil, model.TaskStatusUnknown, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO tasks (public_id, title, description) VALUES ($1, $2, $3)",
		publicID, title, description)
	if err != nil {
		log.Error().Err(err).Msg("failed to create task")
		return nil, model.TaskStatusUnknown, fmt.Errorf("failed to create task: %w", err)
	}

	task := model.Task{}
	err = tx.QueryRow(ctx, "SELECT id, created_at, updated_at FROM tasks WHERE public_id = $1", publicID).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to get task")
		return nil, model.TaskStatusUnknown, fmt.Errorf("failed to get task: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO task_has_status (task_id, status_id, user_id) VALUES ($1, $2, $3)",
		task.ID, model.TaskStatusCreated, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create task status")
		return nil, model.TaskStatusUnknown, fmt.Errorf("failed to create task status: %w", err)
	}

	return &task, model.TaskStatusCreated, nil
}

func (in *Instance) SetTaskStatus(ctx context.Context, publicID uuid.UUID, status model.TaskStatus) error {
	_, err := in.conn.Exec(ctx, "INSERT INTO task_has_status (task_id, status_id) VALUES "+
		"((SELECT id FROM tasks WHERE public_id = $1), $2)", publicID, status)
	if err != nil {
		log.Error().Err(err).Msg("failed to set task status")
		return fmt.Errorf("failed to set task status: %w", err)
	}

	return nil
}
