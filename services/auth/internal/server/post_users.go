package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Nerzal/gocloak/v13"
	dbmodel "github.com/bsergik/tough-dev/services/auth/internal/database/model"
	kafkamodelv1 "github.com/bsergik/tough-dev/services/auth/pkg/message-broker/kafka/model/v1"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostUsers struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Enabled   bool   `json:"enabled"`
}

func (in *Instance) PostUsers(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("cannot read body")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Debug().RawJSON("req_body", reqBody).Msg("PostUsers()")

	reqMsg := &PostUsers{}
	err = json.Unmarshal(reqBody, &reqMsg)
	if err != nil {
		log.Error().Err(err).Msg("cannot unmarshal body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := gocloak.User{
		FirstName: &reqMsg.FirstName,
		LastName:  &reqMsg.LastName,
		Email:     &reqMsg.Email,
		Enabled:   &reqMsg.Enabled,
		Username:  &reqMsg.Username,
	}

	userID, err := in.keycloak.CreateUser(context.TODO(), in.adminToken.AccessToken, in.usersRealm, user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userPublicID := uuid.Must(uuid.NewV4())
	timeNow := time.Now()

	userDb, err := in.db.CreateUser(
		ctx.Request.Context(),
		dbmodel.RoleTypeUser, // TODO replace with the actual role.
		userPublicID,
		*user.Email,
		*user.Username,
		*user.FirstName,
		*user.LastName,
		*user.Enabled,
		timeNow,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = in.mq.PublishUserCreated(context.TODO(), &kafkamodelv1.User{
		PublicId:  userPublicID.String(),
		RoleId:    int32(userDb.RoleID),
		Email:     *user.Email,
		Username:  *user.Username,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Enabled:   *user.Enabled,
		CreatedAt: &timestamppb.Timestamp{Seconds: timeNow.Unix()},
		UpdatedAt: &timestamppb.Timestamp{Seconds: timeNow.Unix()},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to publish user created")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = in.keycloak.SetPassword(context.TODO(), in.adminToken.AccessToken, userID, in.usersRealm, reqMsg.Username, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to set user password")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
