package application

import (
	"context"
	"net/http"

	sloggin "github.com/samber/slog-gin"

	"github.com/gin-gonic/gin"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/service/user"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/password"
)

// SetupRoutes - define route.
func (app *App) SetupRoutes(ctx context.Context) {
	gin.SetMode(app.cfg.GinMode)
	router := gin.New()
	// recovery middleward
	router.Use(sloggin.New(logger.FromContext(ctx)))
	router.Use(gin.Recovery())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})
	app.Router = router
	app.loadUserRoutes()
}

func (app *App) loadUserRoutes() {
	usersGroup := app.Router.Group("/users")
	userStore := user.NewUserStore(app.db)
	usersHandler := user.NewHandler(
		userStore,
		password.NewPasswordHandler(),
	)
	usersHandler.RegisterRoute(usersGroup)
}
