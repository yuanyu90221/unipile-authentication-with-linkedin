package auth

import (
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/config"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/jwt"
)

type Handler struct {
	jwtHandler jwt.JwtHandler
	appConfig  *config.Config
}

func NewHandler(
	jwtHandler jwt.JwtHandler,
	appConfig *config.Config,
) *Handler {
	return &Handler{
		jwtHandler: jwtHandler,
		appConfig:  appConfig,
	}
}
