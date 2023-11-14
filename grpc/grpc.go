package grpc

import (
	"github.com/husanmusa/auth_pro_service/config"
	"github.com/husanmusa/auth_pro_service/genproto/auth_service"
	"github.com/husanmusa/auth_pro_service/grpc/client"
	"github.com/husanmusa/auth_pro_service/grpc/service"
	"github.com/husanmusa/auth_pro_service/storage"
	"github.com/saidamir98/udevs_pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	auth_service.RegisterUserServiceServer(grpcServer, service.NewUserService(cfg, log, strg, svcs))

	reflection.Register(grpcServer)
	return
}
