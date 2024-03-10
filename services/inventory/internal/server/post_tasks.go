package server

import (
	"encoding/json"
	"io"
	"net/http"

	mqmodel "github.com/bsergik/tough-dev/services/inventory/internal/message-broker/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type PostTasks struct {
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (in *Instance) PostTasks(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("cannot read body")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Debug().RawJSON("req_body", reqBody).Msg("PostTasks()")

	reqMsg := &PostTasks{}
	err = json.Unmarshal(reqBody, &reqMsg)
	if err != nil {
		log.Error().Err(err).Msg("cannot unmarshal body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	taskPublicID := uuid.Must(uuid.NewV4())

	task, taskStatus, err := in.db.CreateTask(ctx.Request.Context(), taskPublicID, reqMsg.UserID, reqMsg.Title, reqMsg.Description)
	if err != nil {
		log.Error().Err(err).Msg("failed to create task")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	mqTask := &mqmodel.Task{
		PublicID:    task.PublicID,
		UserID:      reqMsg.UserID,
		Title:       task.Title,
		Description: task.Description,
		StatusID:    int(taskStatus),
	}

	err = in.mq.PublishTaskCreated(ctx.Request.Context(), mqTask)
	if err != nil {
		log.Error().Err(err).Msg("failed to publish task created")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"public_id": task.PublicID,
		"status":    taskStatus,
	})
}
