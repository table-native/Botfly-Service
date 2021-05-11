package service

import (
	"context"

	"github.com/table-native/Botfly-Service/auth"
	pb "github.com/table-native/Botfly-Service/generated"
	"github.com/table-native/Botfly-Service/logger"
	"go.uber.org/zap"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

//
func (u *UserService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	//skip authentication for login service
	return ctx, nil
}

func (u *UserService) Login(ctx context.Context, userId *pb.UserIdentity) (*pb.Token, error) {
	logger.Info("Logging in for ", zap.String("string", userId.EmailId))
	return &pb.Token{Jwt: auth.GetToken(userId.EmailId)}, nil
}
