package main

import (
	"context"
	"fmt"
	"log"
	"net"
	svc "timer-server/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	svc.UnimplementedTimerServer
}

func (s *server) GetTime(ctx context.Context, req *emptypb.Empty) (*svc.GetResponse, error) {
	return &svc.GetResponse{
		Time: timestamppb.Now(),
	}, nil
}

func main() {
	log.Print("starting the service...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen %d: %e", grpcPort, err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	svc.RegisterTimerServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
	log.Printf("server is listening at %v", lis.Addr())
}
