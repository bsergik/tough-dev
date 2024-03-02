package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (in *Instance) GetCerts(ctx *gin.Context) {
	resp, err := in.keycloak.GetCerts(context.TODO(), in.usersRealm)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get certs")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, "text", []byte(resp.String()))
}
