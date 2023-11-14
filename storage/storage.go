package storage

import (
	"context"
	pb "github.com/husanmusa/auth_pro_service/genproto/auth_service"
)

type StorageI interface {
	CloseDB()
	User() UserI
}

type UserI interface {
	CreateUser(ctx context.Context, in *pb.User) error
	UpdateUser(ctx context.Context, in *pb.User) error
	GetUser(ctx context.Context, in *pb.ById) (*pb.User, error)
	GetUserList(ctx context.Context, in *pb.GetUserListRequest) (*pb.GetUserListResponse, error)
	DeleteUser(ctx context.Context, in *pb.ById) error
	GetByUsername(ctx context.Context, username string) (*pb.User, error)
}
