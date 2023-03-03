package rpc

import (
	"github.com/hktalent/proxypool/pkg/rpc/grpc-proxypool"
	"github.com/hktalent/proxypool/pkg/storage"
	"golang.org/x/net/context"
)

type ProxyPoolGRPCService struct{}

func (ProxyPoolGRPCService) Get(ctx context.Context, request *grpc_proxypool.Request) (response *grpc_proxypool.Response, err error) {
	var result string
	switch request.Type {
	case "http":
		result = storage.ProxyRandom().Data
	case "https":
		result = storage.ProxyFind("https").Data
	default:
		result = storage.ProxyRandom().Data
	}
	return &grpc_proxypool.Response{Result: result}, nil
}
