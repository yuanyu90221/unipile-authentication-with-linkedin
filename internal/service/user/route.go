package user

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/config"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/service/auth"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/util"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/jwt"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/password"
)

type Handler struct {
	userStore       *UserStore
	passwordHandler password.PasswordHandler
	authHandler     *auth.Handler
	jwtHandler      jwt.JwtHandler
	appConfig       *config.Config
}

func NewHandler(
	userStore *UserStore,
	passwordHandler password.PasswordHandler,
	authHandler *auth.Handler,
	jwtHandler jwt.JwtHandler,
	appConfig *config.Config,
) *Handler {
	return &Handler{
		userStore:       userStore,
		passwordHandler: passwordHandler,
		authHandler:     authHandler,
		jwtHandler:      jwtHandler,
		appConfig:       appConfig,
	}
}

func (h *Handler) RegisterRoute(router *gin.RouterGroup) {
	router.POST("/", h.CreateUser)
	router.POST("/login", h.Login)
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var request CreateUserRequest
	if err := util.ParseJSON(ctx.Request, &request); err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusBadRequest, err)
		return
	}
	if err := util.Validate.Struct(request); err != nil {
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			util.WriteError(ctx, ctx.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload:%v", valErrs))
		}
		return
	}
	hashedPassword, err := h.passwordHandler.HashPassword(request.Password)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}
	createUserParam := CreateUserParam{Account: request.Account, HashedPassword: hashedPassword}
	createdUser, err := h.userStore.CreateUser(ctx, createUserParam)
	if err != nil {
		if errors.Is(err, ErrorForDuplicateKey) {
			util.WriteError(ctx, ctx.Writer, http.StatusConflict, errors.New("user resource conflict"))
		} else {
			util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		}
		return
	}
	util.FailOnError(
		util.WriteJSON(ctx.Writer, http.StatusCreated, createdUser),
		"failed to create user",
		logger.FromContext(ctx),
	)
}

func (h *Handler) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := util.ParseJSON(ctx.Request, &request); err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusBadRequest, err)
		return
	}
	if err := util.Validate.Struct(request); err != nil {
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			util.WriteError(ctx, ctx.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload:%v", valErrs))
		}
		return
	}
	// find user by account
	user, err := h.userStore.FindByAccount(ctx, request.Account)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, fmt.Errorf("failed to find user %w", err))
		return
	}

	// check hashdedPassword
	if !h.passwordHandler.CheckPassword(request.Password, user.HashedPassword) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "credential incorrect"})
		return
	}
	// gen access token, refresh token
	tokens, err := h.genTokens(user)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	// store refresh token
	_, err = h.userStore.UpdateRefreshToken(ctx, tokens.RefreshToken, user.ID)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("accessToken", tokens.AccessToken)
	ctx.Header("refreshToken", tokens.RefreshToken)
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) genTokens(user UserEntity) (*Tokens, error) {
	accessToken, err := h.jwtHandler.GenerateJWTToken(jwt.JwtSignParam{
		UserID:     user.ID,
		Expiration: time.Second * 3600,
		JwtSecret:  h.appConfig.JWTSecret,
		CurrentTime: func() time.Time {
			return time.Now().UTC()
		},
		Audience: "unipile",
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := h.jwtHandler.GenerateJWTToken(jwt.JwtSignParam{
		UserID:     user.ID,
		Expiration: time.Hour * 3,
		JwtSecret:  h.appConfig.JWTSecret,
		CurrentTime: func() time.Time {
			return time.Now().UTC()
		},
		Audience: "unipile",
	})
	if err != nil {

		return nil, err
	}
	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
