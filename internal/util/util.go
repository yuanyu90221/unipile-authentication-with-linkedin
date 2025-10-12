package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func WriteError(ctx context.Context, w http.ResponseWriter, status int, err error) {
	currentLogger := logger.FromContext(ctx)
	errResp := WriteJSON(w, status, map[string]string{"error": err.Error()})
	if errResp != nil {
		FailOnError(err, "", currentLogger)
	}
}

func FailOnError(err error, msg string, log *slog.Logger) {
	if err != nil && log != nil {
		log.Error(msg, slog.Any("err", err))
		os.Exit(1)
	}
}
