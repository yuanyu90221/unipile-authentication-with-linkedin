package unipile

import (
	"context"
	"fmt"
	"time"

	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/config"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/request"
)

type LinkedinHandler struct {
	requestHandler *request.RequestHandler
	appConfig      *config.Config
}

func NewLinkedinHandler(requestHandler *request.RequestHandler,
	appConfig *config.Config,
) *LinkedinHandler {
	return &LinkedinHandler{
		requestHandler: requestHandler,
		appConfig:      appConfig,
	}
}

func (h LinkedinHandler) ConnectWithCredential(ctx context.Context,
	credentialParam CredentialParam,
) (*ConnectResult, error) {
	var result ConnectResult
	payload := CredentialPayload{
		Provider: "LINKEDIN",
		UserName: credentialParam.UserName,
		Password: credentialParam.Password,
	}
	client := h.requestHandler.Client.SetBaseURL(h.appConfig.UnipileBaseURL)
	resultStatus, err := client.R().
		SetTimeout(6*time.Second).
		SetHeader("X-API-KEY", h.appConfig.UnipileAccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&result).
		Post("/api/v1/accounts")
	if resultStatus.StatusCode() >= 400 {
		return nil, fmt.Errorf("failed to created linked")
	}
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (h LinkedinHandler) ConnectWithCookie(ctx context.Context,
	cookieParam CookieParam,
) (*ConnectResult, error) {
	var result ConnectResult
	payload := CookiePayload{
		Provider:    "LINKEDIN",
		AccessToken: cookieParam.AccessToken,
	}
	client := h.requestHandler.Client.SetBaseURL(h.appConfig.UnipileBaseURL)
	resultStatus, err := client.R().
		SetTimeout(6*time.Second).
		SetHeader("X-API-KEY", h.appConfig.UnipileAccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&result).
		Post("/api/v1/accounts")
	statusCode := resultStatus.StatusCode()
	if statusCode >= 400 {
		return nil, fmt.Errorf("failed to created linked")
	}
	if err != nil {
		return nil, err
	}
	return &result, err
}
