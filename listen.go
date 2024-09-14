package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

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
