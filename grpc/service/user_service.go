package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/husanmusa/auth_pro_service/config"
	pb "github.com/husanmusa/auth_pro_service/genproto/auth_service"
	"github.com/husanmusa/auth_pro_service/grpc/client"
	"github.com/husanmusa/auth_pro_service/pkg/utils"
	"github.com/husanmusa/auth_pro_service/storage"
	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	enforcer *casbin.Enforcer
	services client.ServiceManagerI
	pb.UnimplementedUserServiceServer
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI, enforcer *casbin.Enforcer) *userService {
	return &userService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: svcs,
		enforcer: enforcer,
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

		token := utils.GenerateToken(user.Id, user.Role)
		return &pb.SignInResp{Token: token}, nil
	}

	return nil, status.Error(codes.InvalidArgument, err.Error())
}

func (s *userService) HasAccess(ctx context.Context, req *pb.HasAccessReq) (*pb.HasAccessResp, error) {
	s.log.Info("---HasAccess--->", logger.Any("req", req))
	token, err := utils.ValidateToken(req.Token)
	if err != nil {
		s.log.Error("!!!HasAccess--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.log.Error("!!!HasAccess--->", logger.Error(errors.New("invalid token claims in hasAccess")))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	role := claims["role"]

	err = s.enforcer.LoadPolicy()
	if err != nil {
		s.log.Error("!!!HasAccess--->", logger.Error(errors.New("Failed to load policy from DB")))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Casbin enforces policy
	ok, err = s.enforcer.Enforce(role, req.Obj, req.Act)
	if err != nil {
		s.log.Error("!!!HasAccess--->", logger.Error(errors.New("Error occurred when authorizing user")))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if !ok {
		s.log.Error("!!!HasAccess--->", logger.Error(errors.New("You are not authorized")))
		return nil, status.Error(codes.InvalidArgument, errors.New("You are not authorized").Error())
	}

	return &pb.HasAccessResp{HasAccess: ok}, nil
}
