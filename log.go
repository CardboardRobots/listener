package listener

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func LogGRPC(split bool) grpc.UnaryServerInterceptor {
	if split {
		return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			log.Println(info.FullMethod)
			defer func() {
				if err != nil {
					log.Printf("\terror: %v\n", err)
				}
			}()
			resp, err = handler(ctx, req)
			return
		}
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if err != nil {
				log.Printf("%v\n\terror: %v\n", info.FullMethod, err)
			} else {
				log.Println(info.FullMethod)
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func LogStreamGRPC(split bool) grpc.StreamServerInterceptor {
	if split {
		return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
			log.Println(info.FullMethod)
			defer func() {
				if err != nil {
					log.Printf("\terror: %v\n", err)
				}
			}()
			err = handler(srv, stream)
			return
		}
	}
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if err != nil {
				log.Printf("%v\n\terror: %v\n", info.FullMethod, err)
			} else {
				log.Println(info.FullMethod)
			}
		}()
		err = handler(srv, stream)
		return
	}
}
