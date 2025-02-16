package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/mikhailsoldatkin/book_store/internal/closer"
	"github.com/mikhailsoldatkin/book_store/internal/config"
	"github.com/mikhailsoldatkin/book_store/internal/handlers"
)

// App represents the application with its dependencies and gRPC, HTTP and Swagger servers.
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

// NewApp initializes a new App instance with the given context and sets up the necessary dependencies.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts the HTTP server. It handles graceful shutdown by waiting for context cancellation or termination signals.
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := a.runHTTPServer(); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	cancel()

	wg.Wait()

	return nil
}

// initDeps initializes the dependencies required by the App.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, fun := range inits {
		if err := fun(ctx); err != nil {
			return err
		}
	}

	return nil
}

// initConfig loads the application configuration.
func (a *App) initConfig(_ context.Context) error {
	if _, err := config.Load(); err != nil {
		return err
	}

	return nil
}

// initServiceProvider initializes the service provider used by the application.
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

// initHTTPServer initializes the HTTP server and sets up request handlers.
func (a *App) initHTTPServer(_ context.Context) error {
	db := a.serviceProvider.DBClient()

	muxRouter := initRouter(db)

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.config.HTTP.Address,
		Handler:           muxRouter,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

// initRouter initializes mux.Router
func initRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	for _, r := range handlers.Routes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			HandlerFunc(handlers.WrapHandler(db, r.Handler))
	}

	return router
}

// runHTTPServer starts the HTTP server and listens for incoming requests.
func (a *App) runHTTPServer() error {
	slog.Info(fmt.Sprintf("HTTP server is running on %s", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
