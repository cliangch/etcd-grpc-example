package example

import (
	"golang.org/x/net/context"

	"etcd-grpc-example/server/rpc"
)

type exampleServer struct{}

var ExampleServer exampleServer

func (e exampleServer) Example(ctx context.Context, in *rpc.Request) (*rpc.Response, error) {
	return &rpc.Response{Response: in.Request + ",你好"}, nil
}
