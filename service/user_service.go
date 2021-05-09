package service

import (
	"context"
	"os"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/table-native/Botfly-Service/generated"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (u *UserService) Login(ctx context.Context, userId *pb.UserIdentity) (*pb.Token, error) {
	//TODO: This should be in an env file
	os.Setenv("ACCESS_SECRET", "60de694f-0a61-46f1-8575-5987a24b4abd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = userId.EmailId

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return &pb.Token{}, err
	}
	return &pb.Token{Jwt: token}, nil
}
