package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/jwt"
)

type CtxKey string

var (
	UserIDKey CtxKey = "user_id"
)

func (h *Handler) JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ExtractToken(ctx)
		if len(token) == 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "token not provided"})
			return
		}
		userID, err := h.jwtHandler.VerifyJWTToken(jwt.JwtVerifyParam{
			Token:     token,
			JwtSecret: h.appConfig.JWTSecret,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token not valid",
			})
			return
		}
		ctx.Set(UserIDKey, userID)
		ctx.Next()
	}
}

func ExtractToken(ctx *gin.Context) string {
	bearerToken := ctx.Request.Header.Get("Authorization")
	parseStrings := strings.Split(bearerToken, " ")
	if len(parseStrings) == 2 {
		return parseStrings[1]
	}
	return ""
}

func ExtractUserID(ctx *gin.Context) (int64, error) {
	result, ok := ctx.Get(UserIDKey)
	if !ok {
		return 0, fmt.Errorf("user_id not in ctx")
	}
	userID, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("parse user_id failed")
	}
	return userID, nil
}
