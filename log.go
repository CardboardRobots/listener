package listener

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func LogGRPC(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return
}

func LogStreamGRPC(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	err = handler(srv, stream)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return err
}

func LogRequestGRPC(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println(info.FullMethod)
	return handler(ctx, req)
}

func LogStreamRequestGRPC(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	log.Println(info.FullMethod)
	return handler(srv, stream)
}
