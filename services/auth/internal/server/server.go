package server

import (
	"context"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Instance struct {
	adminToken           *gocloak.JWT
	usersRealm           string
	keycloak             *gocloak.GoCloak
	ginEngine            *gin.Engine
	keycloakClientID     string
	keycloakClientSecret string
}

// type KeyCloak interface {
// 	CreateUser(ctx context.Context, token, realm string, user gocloak.User) (string, error)
// 	LoginAdmin(ctx context.Context, username, password, ream string) (*gocloak.JWT, error)
// 	// LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*gocloak.JWT, error)
// 	SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error
// }

func NewServer(keycloak *gocloak.GoCloak) *Instance {
	keycloakAdminUser, ok := os.LookupEnv("KEYCLOAK_ADMIN_USER")
	if !ok {
		log.Panic().Msg("Env KEYCLOAK_ADMIN_USER is required")
	}

	keycloakAdminPass, ok := os.LookupEnv("KEYCLOAK_ADMIN_PASS")
	if !ok {
		log.Panic().Msg("Env KEYCLOAK_ADMIN_PASS is required")
	}

	keycloakAdminRealm, ok := os.LookupEnv("KEYCLOAK_ADMIN_REALM")
	if !ok {
		log.Panic().Msg("Env KEYCLOAK_ADMIN_REALM is required")
	}

	keycloakClientID, ok := os.LookupEnv("KEYCLOAK_CLIENT_ID")
	if !ok {
		log.Panic().Msg("Env KEYCLOAK_CLIENT_ID is required")
	}

	keycloakClientSecret, ok := os.LookupEnv("KEYCLOAK_CLIENT_SECRET")
	if !ok {
		log.Panic().Msg("Env KEYCLOAK_CLIENT_SECRET is required")
	}

	adminToken, err := keycloak.LoginAdmin(context.TODO(), keycloakAdminUser, keycloakAdminPass, keycloakAdminRealm)
	if err != nil {
		log.Panic().Err(err).Msg("Something wrong with the credentials or url")
	}

	// keycloak.RetrospectToken()

	inst := &Instance{
		keycloak:             keycloak,
		adminToken:           adminToken,
		usersRealm:           keycloakAdminRealm, // TODO replace with realm for users.
		keycloakClientID:     keycloakClientID,
		keycloakClientSecret: keycloakClientSecret,
	}

	ginEngine := gin.Default()
	ginEngine.GET("/", func(ctx *gin.Context) {})

	apiV1Grp := ginEngine.Group("/api/v1")
	apiV1Grp.GET("/certs", inst.GetCerts)
	apiV1Grp.POST("/users/validate", inst.PostUsersValidate)
	apiV1Grp.POST("/users", inst.PostUsers)
	apiV1Grp.POST("/users/login", inst.PostUsersLogin)

	inst.ginEngine = ginEngine

	return inst
}

func (in *Instance) Start(bindAddr string) {
	in.ginEngine.Run(bindAddr)
}
