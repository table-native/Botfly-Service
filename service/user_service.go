package service

import (
	"context"

	"github.com/table-native/Botfly-Service/auth"
	"github.com/table-native/Botfly-Service/db"
	pb "github.com/table-native/Botfly-Service/generated"
	"github.com/table-native/Botfly-Service/logger"
	"go.uber.org/zap"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userDto *db.UserDto
}

func NewUserService(userDto *db.UserDto) *UserService {
	return &UserService{
		userDto: userDto,
	}
}

//
func (u *UserService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	//skip authentication for login service
	return ctx, nil
}

func (u *UserService) Login(ctx context.Context, userId *pb.UserIdentity) (*pb.Token, error) {
	logger.Info("Logging in for ", zap.String("string", userId.EmailId))

	user := <-u.userDto.Create(db.UserModel{
		Email: userId.EmailId,
		Score: 0,
	})

	return &pb.Token{Jwt: auth.GetToken(user.Id)}, nil
}
