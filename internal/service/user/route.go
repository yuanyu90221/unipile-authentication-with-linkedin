package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/util"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/password"
)

type Handler struct {
	userStore       *UserStore
	passwordHandler password.PasswordHandler
}

func NewHandler(
	userStore *UserStore,
	passwordHandler password.PasswordHandler,
) *Handler {
	return &Handler{
		userStore:       userStore,
		passwordHandler: passwordHandler,
	}
}

func (h *Handler) RegisterRoute(router *gin.RouterGroup) {
	router.POST("/", h.CreateUser)
	// router.GET("/:id", h.GetFlightById)
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
