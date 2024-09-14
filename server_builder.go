package listener

import (
	"runtime/debug"
	"slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type ServerBuilder struct {
	unary  []grpc.UnaryServerInterceptor
	stream []grpc.StreamServerInterceptor
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{}
}

func (sb *ServerBuilder) AddUnary(u ...grpc.UnaryServerInterceptor) *ServerBuilder {
	sb.unary = append(sb.unary, u...)
	return sb
}

func (sb *ServerBuilder) AddStream(s ...grpc.StreamServerInterceptor) *ServerBuilder {
	sb.stream = append(sb.stream, s...)
	return sb
}

func (sb *ServerBuilder) Build() *grpc.Server {
	unary, stream := createRecovery()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			append(slices.Clone(sb.unary), unary)...,
		),
		grpc.ChainStreamInterceptor(
			append(slices.Clone(sb.stream), stream)...,
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
