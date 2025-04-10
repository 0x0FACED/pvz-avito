package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/config"
)

type App struct {
	server *http.Server
	config *config.AppConfig
}

func New(srv *http.Server, cfg *config.AppConfig) *App {
	return &App{
		server: srv,
		config: cfg,
	}
}

func (s *App) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

func (s *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
