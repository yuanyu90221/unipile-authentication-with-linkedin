package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtPayload struct {
	UserID int64 `json:"user_id"`
}

type JwtSignParam struct {
	UserID      int64            `json:"user_id"`
	Expiration  time.Duration    `json:"expiration"`
	JwtSecret   string           `json:"jwt_secret"`
	CurrentTime func() time.Time `json:"current_time"`
	Audience    string           `json:"audience"`
}

type JwtVerifyParam struct {
	Token     string `json:"token"`
	JwtSecret string `json:"jwt_secret"`
}

type JwtHandler interface {
	GenerateJWTToken(jwtSignParam JwtSignParam) (string, error)
	VerifyJWTToken(jwtVerifyParam JwtVerifyParam) (int64, error)
}

type Handler struct{}

var _ JwtHandler = (*Handler)(nil)

func NewJwtHandler() JwtHandler {
	return &Handler{}
}

func (h *Handler) GenerateJWTToken(jwtSignParam JwtSignParam) (string, error) {
	claims := jwt.MapClaims{
		"user_id": jwtSignParam.UserID,
		"exp":     jwtSignParam.CurrentTime().UTC().Add(jwtSignParam.Expiration).Unix(),
		"aud":     jwtSignParam.Audience,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSignParam.JwtSecret))
}

func (h *Handler) VerifyJWTToken(jwtVerifyParam JwtVerifyParam) (int64, error) {
	var userID int64
	token, err := jwt.Parse(jwtVerifyParam.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtVerifyParam.JwtSecret), nil
	})
	if err != nil {
		return userID, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, err = strconv.ParseInt(fmt.Sprintf("%.0f", claims["user_id"]), 10, 64)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}
	return 0, nil
}
