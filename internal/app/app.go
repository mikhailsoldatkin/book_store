package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mikhailsoldatkin/book_store/internal/config"
	"github.com/mikhailsoldatkin/book_store/internal/handlers/books"

	"github.com/mikhailsoldatkin/book_store/internal/closer"
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
	wg.Add(1) // add more later ...

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	gracefulShutdown(ctx, cancel, wg)

	return nil
}

// gracefulShutdown handles the termination process by waiting for either a context cancellation
// or termination signals. It cancels the context and waits for all goroutines to finish.
func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		slog.Info("terminating: context cancelled")
	case <-waitSignal():
		slog.Info("terminating: via signal")
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

// waitSignal creates a channel to receive termination signals (SIGINT, SIGTERM).
// It returns the channel to allow waiting for these signals.
func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	return sigterm
}

// initDeps initializes the dependencies required by the App.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initConfig loads the application configuration.
func (a *App) initConfig(_ context.Context) error {
	_, err := config.Load()
	if err != nil {
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
func (a *App) initHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()
	dbClient := a.serviceProvider.DBClient()

	// add handlers
	booksHandler := books.NewHandler(dbClient)

	// register handlers
	mux.HandleFunc("/books", booksHandler.List)

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.config.HTTP.Address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

// runHTTPServer starts the HTTP server and listens for incoming requests.
func (a *App) runHTTPServer() error {
	slog.Info("HTTP server is running on", "address", a.httpServer.Addr)

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
