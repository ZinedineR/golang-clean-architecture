package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"os"
)

func NewGRPCServer(port string) (*grpc.Server, net.Listener) {
	server := grpc.NewServer()
	serverConnection, err := net.Listen("tcp", ":"+port)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return server, serverConnection
}

func NewGRPCClient(server *grpc.Server, net net.Listener, port string) *grpc.ClientConn {
	go func() {
		if err := server.Serve(net); err != nil {
			slog.Error("failed to serve: ", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	client, err := grpc.NewClient(
		"0.0.0.0:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return client
}
