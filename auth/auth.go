package auth

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/table-native/Botfly-Service/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Claims string

var ACCESS_SECRET = "60de694f-0a61-46f1-8575-5987a-24b4abd"

func VerifyToken() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		parsedToken, err := jwt.ParseWithClaims(
			token,
			&jwt.StandardClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(ACCESS_SECRET), nil
			})

		claims, ok := parsedToken.Claims.(*jwt.StandardClaims)

		if !ok || !parsedToken.Valid {
			logger.Fatal("Failed validating token", zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "Bad authorization string")
		}

		newCtx := context.WithValue(ctx, Claims("userId"), claims.Id)
		return newCtx, nil
	}
}

func GetToken(userId string) string {
	//TODO: This should be in an env file
	atClaims := jwt.StandardClaims{}
	atClaims.Id = userId

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(ACCESS_SECRET))
	return token
}
