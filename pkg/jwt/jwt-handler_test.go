package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/jwt"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name         string // description of this test case
		jwtSignParam jwt.JwtSignParam
		want         int64
		wantErr      error
	}{
		{
			name: "test for jwt handler",
			jwtSignParam: jwt.JwtSignParam{
				UserID:     1,
				Expiration: time.Duration(3 * time.Second),
				CurrentTime: func() time.Time {
					return time.Now().UTC()
				},
				JwtSecret: "test_secret",
				Audience:  "test",
			},
			want:    1,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtHandler := jwt.NewJwtHandler()
			token, err := jwtHandler.GenerateJWTToken(tt.jwtSignParam)
			require.NoError(t, err)
			userID, err := jwtHandler.VerifyJWTToken(jwt.JwtVerifyParam{
				Token:     token,
				JwtSecret: tt.jwtSignParam.JwtSecret,
			})
			assert.Equal(t, tt.want, userID)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}
