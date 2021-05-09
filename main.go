package main

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/table-native/Botfly-Service/auth"
	pb "github.com/table-native/Botfly-Service/generated"
	"github.com/table-native/Botfly-Service/logger"
	"github.com/table-native/Botfly-Service/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var port = ":50051"

func main() {
	logger.Info("Starting server at", zap.String("port", port))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logger.Get()),
			grpc_auth.UnaryServerInterceptor(auth.VerifyToken()),
		)),
	)
	pb.RegisterUserServiceServer(s, &service.UserService{})
	pb.RegisterGameServiceServer(s, &service.GameService{})

	if err := s.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}
