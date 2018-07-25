package discover

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var client *clientv3.Client

// 服务注册
func Register(endpoints []string, name, addr string, ttl int64) error {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 10 * time.Second,
	})

	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second * time.Duration(ttl))
	go func() {
		for {
			response, err := client.Get(context.Background(), fmt.Sprintf("/%s/%s/%s", schema, name, addr))
			if err != nil {
				log.Printf("endpoints:%v, name:%s, addr:%s, error:%v", endpoints, name, addr, err)
			} else if response.Count == 0 {
				keepAlive(name, addr, ttl)
			}

			<-ticker.C
		}
	}()

	return err
}

func keepAlive(name, addr string, ttl int64) error {
	response, err := client.Grant(context.Background(), ttl)

	if err != nil {
		return err
	}

	if _, err = client.Put(context.Background(), fmt.Sprintf("/%s/%s/%s", schema, name, addr), addr, clientv3.WithLease(response.ID)); err != nil {
		return err
	}
	_, err = client.KeepAlive(context.Background(), response.ID)

	return err
}

// remove server from etcd
func Remove(name, addr string) {
	if client == nil {
		return
	}
	client.Delete(context.Background(), fmt.Sprintf("/%s/%s/%s", schema, name, addr))
}
