package server

import (
	"context"

	dbmodel "github.com/bsergik/tough-dev/services/inventory/internal/database/model"
	mqmodel "github.com/bsergik/tough-dev/services/inventory/internal/message-broker/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type Instance struct {
	ginEngine *gin.Engine
	db        Database
	mq        MessageQueue
}

type Database interface {
	CreateTask(ctx context.Context, publicID, userID uuid.UUID, title, description string) (*dbmodel.Task, dbmodel.TaskStatus, error)
}

type MessageQueue interface {
	PublishTaskCreated(ctx context.Context, task *mqmodel.Task) error
}

func NewServer(db Database, mq MessageQueue) *Instance {
	inst := &Instance{
		db: db,
		mq: mq,
	}

	ginEngine := gin.Default()
	ginEngine.GET("/", func(ctx *gin.Context) {})

	apiV1Grp := ginEngine.Group("/api/v1")
	apiV1Grp.POST("/tasks", inst.PostTasks)

	inst.ginEngine = ginEngine

	return inst
}

func (in *Instance) Start(bindAddr string) {
	in.ginEngine.Run(bindAddr)
}
