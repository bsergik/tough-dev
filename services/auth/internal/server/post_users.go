package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
		log.Error().Stack().Err(err).Msg("cannot read body")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Debug().RawJSON("req_body", reqBody).Msg("PostUsers()")

	reqMsg := &PostUsers{}
	err = json.Unmarshal(reqBody, &reqMsg)
	if err != nil {
		log.Error().Stack().Err(err).Msg("cannot unmarshal body")
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
		log.Error().Stack().Err(err).Msg("failed to create user")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = in.keycloak.SetPassword(context.TODO(), in.adminToken.AccessToken, userID, in.usersRealm, reqMsg.Username, true)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to set user password")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
