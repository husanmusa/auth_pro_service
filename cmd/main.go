package main

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/husanmusa/auth_pro_service/api"
	"github.com/husanmusa/auth_pro_service/api/handlers"
	"github.com/husanmusa/auth_pro_service/config"
	"github.com/husanmusa/auth_pro_service/grpc"
	"github.com/husanmusa/auth_pro_service/grpc/client"
	"github.com/husanmusa/auth_pro_service/storage/postgres"
	"github.com/saidamir98/udevs_pkg/logger"
	"net"
)

func main() {
	var loggerLevel string
	cfg := config.Load()

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
	case config.TestMode:
		loggerLevel = logger.LevelDebug
	default:
		loggerLevel = logger.LevelInfo
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer func() {
		if err := logger.Cleanup(log); err != nil {
			log.Error("Failed to cleanup logger", logger.Error(err))
		}
	}()

	pgStore, err := postgres.NewPostgres(context.Background(), cfg, log)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	defer pgStore.CloseDB()

	svcs, err := client.NewGrpcClients(cfg)
	if err != nil {
		log.Panic("client.NewGrpcClients", logger.Error(err))
	}
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", "config/rbac_policy.csv")
	if err != nil {
		log.Error("casbin.NewEnforcer", logger.Error(err))
	}
	grpcServer := grpc.SetUpServer(cfg, log, pgStore, svcs, enforcer)
	go func() {
		lis, err := net.Listen("tcp", cfg.AuthServicePort)
		if err != nil {
			log.Error("net.Listen", logger.Error(err))
		}

		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpcServer.Serve", logger.Error(err))
		}
	}()
	h := handlers.NewHandler(cfg, log, svcs)


	r := api.SetUpRouter(h, enforcer)

	err = r.Listen(cfg.HTTPPort)
	if err != nil {
		panic(err)
	}

	log.Info("server shutdown successfully")
}
