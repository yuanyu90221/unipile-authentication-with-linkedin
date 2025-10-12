package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/config"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/db"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/util"
)

// App define app dependency.
type App struct {
	Router *gin.Engine
	cfg    *config.Config
	db     *sql.DB
}

// New - app constructor.
func New(ctx context.Context, cfg *config.Config) *App {
	log := logger.FromContext(ctx)
	// setup DB connection
	dbConn, err := db.Connect(cfg.DBURL)
	if err != nil {
		util.FailOnError(err, "failed to connect", log)
	}
	app := &App{
		cfg: cfg,
		db:  dbConn,
	}
	// setup routes
	app.SetupRoutes(ctx)
	return app
}

// Start - app 啟動.
func (app *App) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", app.cfg.Port),
		Handler:           app.Router,
		ReadHeaderTimeout: time.Minute,
	}
	// graceful shutdown close redis
	defer func() {
		if err := app.db.Close(); err != nil {
			// log.Println("failed to close db connection", err)
			util.FailOnError(err, "failed to close db connection", log)
		}
	}()
	log.Info(fmt.Sprintf("Starting server on %s", app.cfg.Port))
	errCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("server cancel")
		timeout, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		err := server.Shutdown(timeout)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	}
}
