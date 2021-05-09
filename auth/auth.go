package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Claims string

func VerifyToken() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			//Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Bad authorization string")
		}

		newCtx := context.WithValue(ctx, Claims("userId"), claims)
		return newCtx, nil
	}
}

func GetToken(userId string) string {
	//TODO: This should be in an env file
	os.Setenv("ACCESS_SECRET", "60de694f-0a61-46f1-8575-5987a24b4abd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = userId

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	return token
}
