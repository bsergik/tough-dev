package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type PostUsersValidate struct {
	Token string `json:"token"`
}

func (in *Instance) PostUsersValidate(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("cannot read body")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	reqMsg := &PostUsersValidate{}
	err = json.Unmarshal(reqBody, &reqMsg)
	if err != nil {
		log.Error().Stack().Err(err).Msg("cannot unmarshal body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	jwt, err := in.keycloak.RetrospectToken(
		context.TODO(),
		reqMsg.Token,
		in.keycloakClientID,
		in.keycloakClientSecret,
		in.usersRealm,
	)
	if err != nil {
		log.Error().Stack().Err(err).Str("token", reqMsg.Token).Msg("failed to validate token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// if jwt.Active != nil && !*jwt.Active {
	// 	jwt, err := in.keycloak.RefreshToken(
	// 		context.TODO(),
	// 		reqMsg.Token,
	// 		in.keycloakClientID,
	// 		in.keycloakClientSecret,
	// 		in.usersRealm,
	// 	)
	// 	if err != nil {
	// 		log.Error().Stack().Err(err).Str("token", ).Msg("failed to validate token")
	// 		ctx.AbortWithStatus(http.StatusUnauthorized)
	// 		return
	// 	}
	// }

	// jwt.Active

	ctx.Data(http.StatusOK, "text", []byte(jwt.String()))
}
