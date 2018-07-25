package main

import (
	"net"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"etcd-grpc-example/server/rpc"
	"etcd-grpc-example/discover"
	"etcd-grpc-example/server/example"
)

func main() {
	addr := "127.0.0.1:10086"

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net listener error:%v", err)
	}

	s := grpc.NewServer()
	rpc.RegisterExampleServer(s, example.ExampleServer)

	endpoints := []string{"127.0.0.1:2379"}
	name := "example-server"

	if err := discover.Register(endpoints, name, addr, 10); err != nil {
		log.Fatalf("register etcd error:%v", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-ch
		discover.Remove(name, addr)
		os.Exit(1)
	}()

	log.Printf("server start addr:%s", addr)
	s.Serve(listener)
}
