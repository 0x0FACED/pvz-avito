package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	auth_svc "github.com/0x0FACED/pvz-avito/internal/auth/application"
	"github.com/0x0FACED/pvz-avito/internal/auth/delivery/http"
	auth_db "github.com/0x0FACED/pvz-avito/internal/auth/infra/postgres"
	"github.com/0x0FACED/pvz-avito/internal/pkg/config"
	"github.com/0x0FACED/pvz-avito/internal/pkg/database"
	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
	product_svc "github.com/0x0FACED/pvz-avito/internal/product/application"
	product_db "github.com/0x0FACED/pvz-avito/internal/product/infra/postgres"
	pvz_svc "github.com/0x0FACED/pvz-avito/internal/pvz/application"
	pvz_db "github.com/0x0FACED/pvz-avito/internal/pvz/infra/postgres"
	reception_svc "github.com/0x0FACED/pvz-avito/internal/reception/application"
	reception_db "github.com/0x0FACED/pvz-avito/internal/reception/infra/postgres"
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
		log.Fatalln("cant create logger, err: ", err)
		return
	}

	// init all loggers with features
	httpLogger := logger.WithFeature("http")
	authSvcLogger := logger.WithFeature("auth_svc")
	pvzSvcLogger := logger.WithFeature("pvz_svc")
	productSvcLogger := logger.WithFeature("product_svc")
	receptionSvcLogger := logger.WithFeature("reception_svc")

	// connect to db pool
	pool, err := database.ConnectPool(ctx, cfg.Database)

	// creating all repos
	authRepo := auth_db.NewAuthPostgresRepository(pool)
	pvzRepo := pvz_db.NewPVZPostgresRepository(pool)
	productRepo := product_db.NewProductPostgresRepository(pool)
	receptionRepo := reception_db.NewReceptionPostgresRepository(pool)

	// creating all svcs
	authSvc := auth_svc.NewAuthService(authRepo)
	pvzSvc := pvz_svc.NewPVZService(pvzRepo)
	productSvc := product_svc.NewProductService(productRepo)
	receptionSvc := reception_svc.NewReceptionService(receptionRepo)

	// jwt manager (move diration to cfg)
	jwt := httpcommon.NewManager(cfg.Server.JWTSecret, 256*time.Hour)
	// create all handlers
	authHandler := http.NewHandler(authSvc, jwt)
}
