package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/app"
	auth_svc "github.com/0x0FACED/pvz-avito/internal/auth/application"
	auth_http "github.com/0x0FACED/pvz-avito/internal/auth/delivery/http"
	auth_db "github.com/0x0FACED/pvz-avito/internal/auth/infra/postgres"
	"github.com/0x0FACED/pvz-avito/internal/pkg/config"
	"github.com/0x0FACED/pvz-avito/internal/pkg/database"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	"github.com/0x0FACED/pvz-avito/internal/pkg/middleware"
	product_svc "github.com/0x0FACED/pvz-avito/internal/product/application"
	product_http "github.com/0x0FACED/pvz-avito/internal/product/delivery/http"
	product_db "github.com/0x0FACED/pvz-avito/internal/product/infra/postgres"
	pb "github.com/0x0FACED/pvz-avito/internal/pvz/delivery/grpc/v1"

	pvz_svc "github.com/0x0FACED/pvz-avito/internal/pvz/application"
	pvz_grpc "github.com/0x0FACED/pvz-avito/internal/pvz/delivery/grpc"
	pvz_http "github.com/0x0FACED/pvz-avito/internal/pvz/delivery/http"
	pvz_db "github.com/0x0FACED/pvz-avito/internal/pvz/infra/postgres"
	reception_svc "github.com/0x0FACED/pvz-avito/internal/reception/application"
	reception_http "github.com/0x0FACED/pvz-avito/internal/reception/delivery/http"
	reception_db "github.com/0x0FACED/pvz-avito/internal/reception/infra/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error recovered: %v\nStack trace:\n%s", err, debug.Stack())
		}
	}()

	logger, err := logger.NewZerologLogger(cfg.Logger)
	if err != nil {
		log.Panicln("cant create logger, err: ", err)
		return
	}

	// init all loggers with features
	appLogger := logger.WithFeature("app")
	httpLogger := logger.WithFeature("http")
	authSvcLogger := logger.WithFeature("auth_svc")
	pvzSvcLogger := logger.WithFeature("pvz_svc")
	productSvcLogger := logger.WithFeature("product_svc")
	receptionSvcLogger := logger.WithFeature("reception_svc")

	appLogger.Info().Msg("Loggers with features created")

	appLogger.Info().Msg("Connecting to database...")
	// connect to db pool
	pool, err := database.ConnectPool(ctx, cfg.Database)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	appLogger.Info().Msg("Successfully connected to database")

	// creating all repos
	authRepo := auth_db.NewAuthPostgresRepository(pool)
	pvzRepo := pvz_db.NewPVZPostgresRepository(pool)
	productRepo := product_db.NewProductPostgresRepository(pool)
	receptionRepo := reception_db.NewReceptionPostgresRepository(pool)

	appLogger.Info().Msg("Repos for application services created")

	// creating all svcs
	authSvc := auth_svc.NewAuthService(authRepo, authSvcLogger)
	pvzSvc := pvz_svc.NewPVZService(pvzRepo, receptionRepo, productRepo, pvzSvcLogger)
	productSvc := product_svc.NewProductService(productRepo, receptionRepo, productSvcLogger)
	receptionSvc := reception_svc.NewReceptionService(receptionRepo, receptionSvcLogger)

	appLogger.Info().Msg("Application services created")

	// jwt manager (move diration to cfg)
	jwt := httpcommon.NewManager(cfg.Server.JWTSecret, time.Hour*240)

	appLogger.Info().Msg("JWT Manager created")

	// create middleware
	middleware := middleware.NewMiddlewareHandler(jwt, httpLogger)

	appLogger.Info().Msg("Middleware instance created")

	// create all handlers
	authHandler := auth_http.NewHandler(authSvc, jwt)
	pvzHandler := pvz_http.NewHandler(pvzSvc)
	productHandler := product_http.NewHandler(productSvc)
	receptionHandler := reception_http.NewHandler(receptionSvc)

	appLogger.Info().Msg("Handlers created")

	// registering routes with middleware
	mux := http.NewServeMux()
	authHandler.RegisterRoutes(mux)

	// protected with auth middleware
	privateMux := http.NewServeMux()
	pvzHandler.RegisterRoutes(privateMux)
	productHandler.RegisterRoutes(privateMux)
	receptionHandler.RegisterRoutes(privateMux)

	// apply auth for '/' routes (all expect public /login, /dummyLogin, /register)
	mux.Handle("/", middleware.Auth(privateMux))

	// apply logger middleware for all routes
	// this is final mux
	loggedMux := middleware.Logger(mux)

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	appLogger.Info().Msg("All routes added")

	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      loggedMux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	appLogger.Info().Msg("Application server created")

	var metricsSrv *http.Server

	if cfg.Metrics.Enabled {
		metricsSrv = &http.Server{
			Addr:         cfg.Server.Host + ":" + cfg.Metrics.Port,
			Handler:      metricsMux,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		}

		appLogger.Info().Msg("Metrics server created")
	}

	// adding grpc server
	grpcServer := grpc.NewServer()
	pvzGrpcHandler := pvz_grpc.NewGRPCHandler(pvzSvc)
	pb.RegisterPVZServiceServer(grpcServer, pvzGrpcHandler)

	app := app.New(srv, metricsSrv, grpcServer, appLogger, cfg)

	appLogger.Info().Msg("App instance created, starting servers...")

	go func() {
		if err := app.Start(ctx); err != nil {
			return
		}
	}()

	<-ctx.Done()

	if err := app.Shutdown(); err != nil {
		return
	}
}
