package main

import (
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"github.com/table-native/Botfly-Service/auth"
	"github.com/table-native/Botfly-Service/db"
	pb "github.com/table-native/Botfly-Service/generated"
	"github.com/table-native/Botfly-Service/logger"
	"github.com/table-native/Botfly-Service/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var port = ":50051"
var webPort = ":8080"

func StartGrpcServer() *grpc.Server {
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

	botflyDb := db.NewBotflyDb()
	userDto := db.NewUserDto(botflyDb)
	scriptsDto := db.NewScriptsDto(botflyDb)

	pb.RegisterUserServiceServer(s, service.NewUserService(userDto))
	pb.RegisterGameServiceServer(s, service.NewGameService(scriptsDto))

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()
	return s
}

func buildServer(wrappedGrpc *grpcweb.WrappedGrpcServer) *http.Server {
	return &http.Server{
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			wrappedGrpc.ServeHTTP(resp, req)
		}),
	}
}

func main() {
	grpcServer := StartGrpcServer()
	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}))

	webServer := buildServer(wrappedGrpc)
	http.Handle("/metrics", promhttp.Handler())
	webListener, err := net.Listen("tcp", webPort)

	if err != nil {
		logger.Fatal("Failed starting web server", zap.Error(err))
	}
	logger.Info("Starting web server at ", zap.String("port", webPort))

	if err := webServer.Serve(webListener); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}
