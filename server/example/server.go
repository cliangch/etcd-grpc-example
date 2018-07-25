package example

import (
	"etcd-grpc-example/server/rpc"

	"golang.org/x/net/context"
)

type exampleServer struct{}

var ExampleServer exampleServer

func (e exampleServer) Example(ctx context.Context, in *rpc.Request) (*rpc.Response, error) {
	return &rpc.Response{Response: in.Request + ",你好"}, nil
}
