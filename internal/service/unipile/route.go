package unipile

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/service/auth"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/util"
)

type Handler struct {
	unipileStore    *UnipileStore
	linkedinHandler *LinkedinHandler
	authHandler     *auth.Handler
}

func NewHandler(
	unipileStore *UnipileStore,
	linkedinHandler *LinkedinHandler,
	authHandler *auth.Handler,
) *Handler {
	return &Handler{
		unipileStore:    unipileStore,
		linkedinHandler: linkedinHandler,
		authHandler:     authHandler,
	}
}

func (h *Handler) RegisterRoute(router *gin.RouterGroup) {
	router.Use(h.authHandler.JwtAuthMiddleware())
	router.POST("/credential", h.ConnectUserWithCredential)
	router.POST("/cookie", h.ConnectUserWithCookie)
	router.GET("/", h.ListFederaByUserID)
}

// ConnectUserWithCredential -  linked user with linkedin credential handler
func (h *Handler) ConnectUserWithCredential(ctx *gin.Context) {
	userID, err := auth.ExtractUserID(ctx)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusBadRequest,
			err,
		)
		return
	}
	var request ConnectUserWithCredentialRequest
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

	connectResult, err := h.linkedinHandler.ConnectWithCredential(ctx, CredentialParam{
		UserName: request.Account,
		Password: request.Password,
	})
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	status := connectResult.Object
	if connectResult.CheckPoint.Type != "" {
		status = fmt.Sprintf("%s-%s", status, connectResult.CheckPoint.Type)
	}

	param := CreateUnipileUserFederaParam{
		AccountID: connectResult.AccountID,
		Provider:  "LINKEDIN",
		UserID:    userID,
		Status:    status,
	}
	linkedResult, err := h.unipileStore.CreateUnipileUserFederal(ctx, param)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}
	util.FailOnError(
		util.WriteJSON(ctx.Writer, http.StatusCreated, linkedResult),
		"failed to linked user with unipile credential",
		logger.FromContext(ctx),
	)
}

// ConnectUserWithCookie - linked user with linkedin cookie handler
func (h *Handler) ConnectUserWithCookie(ctx *gin.Context) {
	userID, err := auth.ExtractUserID(ctx)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusBadRequest,
			err,
		)
		return
	}
	var request ConnectUserWithCookieRequest
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

	connectResult, err := h.linkedinHandler.ConnectWithCookie(ctx, CookieParam{
		AccessToken: request.AccessToken,
	})

	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	status := connectResult.Object
	if connectResult.CheckPoint.Type != "" {
		status = fmt.Sprintf("%s-%s", status, connectResult.CheckPoint.Type)
	}
	param := CreateUnipileUserFederaParam{
		AccountID: connectResult.AccountID,
		Provider:  "LINKEDIN",
		UserID:    userID,
		Status:    status,
	}
	linkedResult, err := h.unipileStore.CreateUnipileUserFederal(ctx, param)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}
	util.FailOnError(
		util.WriteJSON(ctx.Writer, http.StatusCreated, linkedResult),
		"failed to linked user with unipile cookie",
		logger.FromContext(ctx),
	)
}

func (h *Handler) ListFederaByUserID(ctx *gin.Context) {
	userID, err := auth.ExtractUserID(ctx)
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusBadRequest,
			err,
		)
		return
	}
	result, err := h.unipileStore.ListUnipileUserFederalByUserID(ctx, ListFederaParam{
		UserID: userID,
	})
	if err != nil {
		util.WriteError(ctx, ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	util.FailOnError(
		util.WriteJSON(ctx.Writer, http.StatusOK, result),
		fmt.Sprintf("failed to list federal with user id :%v", userID),
		logger.FromContext(ctx),
	)
}
