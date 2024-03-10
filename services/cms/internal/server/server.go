package server

import (
	"github.com/bsergik/tough-dev/services/cms/internal/database/psql"
	"github.com/bsergik/tough-dev/services/cms/internal/message-broker/kafka"
	"github.com/gin-gonic/gin"
)

type Instance struct {
	ginEngine *gin.Engine
	db        *psql.Instance
	mq        *kafka.Instance
}

func NewServer(db *psql.Instance, mq *kafka.Instance) *Instance {
	inst := &Instance{
		ginEngine: gin.Default(),
		db:        db,
		mq:        mq,
	}

	inst.ginEngine.GET("/", func(ctx *gin.Context) {})

	// apiV1Grp := inst.ginEngine.Group("/api/v1")

	return inst
}

func (in *Instance) Start(bindAddr string) {
	in.ginEngine.Run(bindAddr)
}
