package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func CreateServer(interceptors ...grpc.UnaryServerInterceptor) *grpc.Server {
	unary, stream := createRecovery()

	chain := make([]grpc.UnaryServerInterceptor, 0, len(interceptors)+1)
	chain = append(chain, unary)
	chain = append(chain, interceptors...)

	// You can now create a server with logging instrumentation that e.g. logs when the unary or stream call is started or finished.
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			chain...,
		),
		grpc.ChainStreamInterceptor(
			stream,
		),
	)

	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(server, healthcheck)
	reflection.Register(server)

	return server
}

func createRecovery() (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) error {
			debug.PrintStack()
			return status.Errorf(codes.Unknown, "panic triggered: %v", p)
		}),
	}
	unary := recovery.UnaryServerInterceptor(opts...)
	stream := recovery.StreamServerInterceptor(opts...)
	return unary, stream
}

func LogGRPC(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return
}

func ErrorInterceptor(errorMap func(err error) codes.Code) func(context.Context, any, *grpc.UnaryServerInfo, grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			code := errorMap(err)
			err = status.Error(code, err.Error())
		}
		return
	}
}

func Listen(ctx context.Context, port0, port1 int, handler http.Handler, server *grpc.Server) {
	listener0, err := getListener(port0)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	listener1, err := getListener(port1)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	go ListenGRPC(cancel, listener0, server)
	go ListenHttp(cancel, listener1, handler)

	<-ctx.Done()
}

func getListener(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

func ListenHttp(cancel context.CancelFunc, listener net.Listener, handler http.Handler) {
	log.Printf("http listening at %v", listener.Addr())
	if err := http.Serve(listener, handler); err != nil {
		log.Fatalf("failed to serve: %v", err)
		cancel()
	}

	cancel()
}

func ListenGRPC(cancel context.CancelFunc, listener net.Listener, server *grpc.Server) {
	log.Printf("gRPC listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
		cancel()
	}

	cancel()
}
