package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/config"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	server        *http.Server
	metricsServer *http.Server
	grpcServer    *grpc.Server

	log    *logger.ZerologLogger
	config *config.AppConfig
}

func New(srv *http.Server, merticsSrv *http.Server, grpcSrv *grpc.Server, l *logger.ZerologLogger, cfg *config.AppConfig) *App {
	return &App{
		server:        srv,
		metricsServer: merticsSrv,
		grpcServer:    grpcSrv,
		log:           l,
		config:        cfg,
	}
}

func (a *App) Start(ctx context.Context) error {
	errChan := make(chan error, 3)

	go func() {
		a.log.Info().Str("address", a.server.Addr).Msg("Starting application server")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	if a.config.GRPCPVZ.Enabled && a.grpcServer != nil {
		go func() {
			addr := ":" + a.config.GRPCPVZ.Port
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				errChan <- err
				return
			}

			a.log.Info().Str("address", addr).Msg("Starting gRPC server")
			if err := a.grpcServer.Serve(lis); err != nil {
				errChan <- err
			}
		}()
	}

	if a.config.Metrics.Enabled && a.metricsServer != nil {
		go func() {
			a.log.Info().Str("address", a.metricsServer.Addr).Msg("Starting metrics server")
			if err := a.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.log.Info().Msg("Shutting down servers...")

	var retErr error

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error().Err(err).Msg("Failed to shutdown application server")
		retErr = err
	} else {
		a.log.Info().Msg("Application server stopped")
	}

	if a.config.GRPCPVZ.Enabled && a.grpcServer != nil {
		a.grpcServer.GracefulStop()
		a.log.Info().Msg("gRPC server stopped")
	}

	if a.config.Metrics.Enabled && a.metricsServer != nil {
		if err := a.metricsServer.Shutdown(ctx); err != nil {
			a.log.Error().Err(err).Msg("Failed to shutdown metrics server")
			retErr = err
		} else {
			a.log.Info().Msg("Metrics server stopped")
		}
	}

	return retErr
}
