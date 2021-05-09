package main

import (
	"log"
	"net"

	pb "github.com/table-native/Botfly-Service/generated"
	"github.com/table-native/Botfly-Service/service"
	"google.golang.org/grpc"
)

var port = ":50051"

func main() {
	log.Printf("Starting server at %s\n", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &service.UserService{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
