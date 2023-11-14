package service

import (
	"context"
	"fmt"
	"github.com/husanmusa/auth_pro_service/config"
	pb "github.com/husanmusa/auth_pro_service/genproto/auth_service"
	"github.com/husanmusa/auth_pro_service/grpc/client"
	"github.com/husanmusa/auth_pro_service/storage"
	"github.com/husanmusa/auth_pro_service/utils"
	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	pb.UnimplementedUserServiceServer
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) *userService {
	return &userService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: svcs,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *pb.User) (*emptypb.Empty, error) {
	s.log.Info("---CreateUser--->", logger.Any("req", req))
	err := s.strg.User().CreateUser(ctx, req)
	if err != nil {
		s.log.Error("!!!For Message--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.log.Info("---CreateUser--->", logger.Any("resp", "SUCCESS in Service"))

	return &emptypb.Empty{}, err
}

func (s *userService) GetUser(ctx context.Context, req *pb.ById) (*pb.User, error) {
	resp, err := s.strg.User().GetUser(ctx, req)

	if err != nil {
		s.log.Error("!!!For Message--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *userService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	s.log.Info("---GetUserList--->", logger.Any("req", req))

	res, err := s.strg.User().GetUserList(ctx, req)

	if err != nil {
		s.log.Error("!!!GetUserList--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, err
}

func (s *userService) UpdateUser(ctx context.Context, req *pb.User) (*emptypb.Empty, error) {
	s.log.Info("---UpdateUser--->", logger.Any("req", req))

	err := s.strg.User().UpdateUser(ctx, req)
	if err != nil {
		s.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &emptypb.Empty{}, err
}

func (s *userService) DeleteUser(ctx context.Context, req *pb.ById) (*emptypb.Empty, error) {
	s.log.Info("---DeleteUser--->", logger.Any("req", req))

	err := s.strg.User().DeleteUser(ctx, req)

	if err != nil {
		s.log.Error("!!!DeleteUser--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *userService) SignInUser(ctx context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	s.log.Info("---GetByUsername--->", logger.Any("req", req))

	user, err := s.strg.User().GetByUsername(ctx, req.Username)
	if err != nil {
		s.log.Error("!!!GetByUsername--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if isTrue := utils.ComparePassword(user.Password, req.Password); isTrue {
		fmt.Println("user before", user.Id)
		//dbUserID, err := strconv.ParseUint(user.Id, 10, 64)
		//if err != nil {
		//	s.log.Error("!!!GetByUsername--->", logger.Error(err))
		//	return nil, status.Error(codes.InvalidArgument, err.Error())
		//}

		token := utils.GenerateToken((user.Id))
		return &pb.SignInResp{Token: token}, nil
	}

	return nil, status.Error(codes.InvalidArgument, err.Error())
}
