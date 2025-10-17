package application

import (
	"context"
	"embed"
	"net/http"

	"github.com/gin-contrib/static"
	sloggin "github.com/samber/slog-gin"

	"github.com/gin-gonic/gin"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/service/auth"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/service/unipile"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/service/user"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/jwt"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/password"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/pkg/request"
)

var ServerFs embed.FS

// SetupRoutes - define route.
func (app *App) SetupRoutes(ctx context.Context) {
	gin.SetMode(app.cfg.GinMode)
	router := gin.New()

	// recovery middleward
	router.Use(sloggin.New(logger.FromContext(ctx)))
	router.Use(gin.Recovery())

	// setup static
	fs, err := static.EmbedFolder(ServerFs, "static")
	if err != nil {
		panic(err)
	}
	router.Use(static.Serve("/", fs))

	// for history mode
	router.NoRoute(func(c *gin.Context) {
		// c.HTML(http.StatusOK, "index.html", nil)
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	app.Router = router
	jwtHandler := jwt.NewJwtHandler()
	authHandler := auth.NewHandler(jwtHandler, app.cfg)
	app.loadUserRoutes(authHandler, jwtHandler)
	app.loadUnipileRoutes(authHandler)
}

func (app *App) loadUserRoutes(authHandler *auth.Handler,
	jwtHandler jwt.JwtHandler,
) {
	usersGroup := app.Router.Group("/users")
	userStore := user.NewUserStore(app.db)
	usersHandler := user.NewHandler(
		userStore,
		password.NewPasswordHandler(),
		authHandler,
		jwtHandler,
		app.cfg,
	)
	usersHandler.RegisterRoute(usersGroup)
}

func (app *App) loadUnipileRoutes(authHandler *auth.Handler) {
	unipileGroup := app.Router.Group("/unipile")
	unipileStore := unipile.NewUnipileStore(app.db)
	linkedInHandler := unipile.NewLinkedinHandler(
		request.NewRequestHandler(),
		app.cfg,
	)
	unipileHandler := unipile.NewHandler(
		unipileStore,
		linkedInHandler,
		authHandler)
	unipileHandler.RegisterRoute(unipileGroup)
}
