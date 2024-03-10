package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type PostUsersLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (in *Instance) PostUsersLogin(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("cannot read body")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	reqMsg := &PostUsersLogin{}
	err = json.Unmarshal(reqBody, &reqMsg)
	if err != nil {
		log.Error().Stack().Err(err).Msg("cannot unmarshal body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	jwt, err := in.keycloak.Login(
		context.TODO(),
		in.keycloakClientID,
		in.keycloakClientSecret,
		in.usersRealm,
		reqMsg.Username,
		reqMsg.Password,
	)
	if err != nil {
		log.Error().Stack().Err(err).Str("username", reqMsg.Username).Msg("failed to login user")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Data(http.StatusOK, "text", []byte(jwt.AccessToken))
}
