package listener

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

type Metadata map[string][]string

const (
	authorizationHeader = "authorization"
)

func GetMetadataGRPC(ctx context.Context) Metadata {
	m, _ := metadata.FromIncomingContext(ctx)
	return Metadata(m)
}

func GetMetadataHTTP(r *http.Request) Metadata {
	return Metadata(r.Header)
}

func (md Metadata) Get(key string) []string {
	return md[key]
}

func (md Metadata) GetFirst(key string) string {
	data := md[key]
	if len(data) == 0 {
		return ""
	}
	return data[0]
}

func (md Metadata) Authorization() string {
	return md.GetFirst(authorizationHeader)
}
