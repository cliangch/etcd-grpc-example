package client

import (
	"fmt"
	"log"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"

	"etcd-grpc-example/discover"
	"etcd-grpc-example/client/rpc"
	"golang.org/x/net/context"
)

var (
	schemaTarget = "super:///%s"
	server       = "example-server"
)

func Example() string {
	endpoints := []string{"127.0.0.1:2379"}

	r := discover.NewResolver(endpoints)
	resolver.Register(r)

	if conn, err := grpc.Dial(fmt.Sprintf(schemaTarget, server), grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name)); err != nil {
		log.Printf("conn error:%v", err)
		return ""
	} else {
		defer conn.Close()

		c := rpc.NewExampleClient(conn)
		if response, err := c.Example(context.Background(), &rpc.Request{Request: "Hi"}); err != nil {
			log.Printf("request error:%v", err)
			return ""
		} else {
			return response.Response
		}
	}
}
