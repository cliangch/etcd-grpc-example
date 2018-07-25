package discover

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
)

const schema = "super"

type etcdResolver struct {
	endpoints []string
	cc        resolver.ClientConn
	client    *clientv3.Client
}

func NewResolver(endpoints []string) resolver.Builder {
	return &etcdResolver{endpoints: endpoints}
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var err error

	if r.client == nil {
		r.client, err = clientv3.New(clientv3.Config{
			Endpoints:   r.endpoints,
			DialTimeout: 10 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}

	r.cc = cc
	go r.watch(fmt.Sprintf("/%s/%s/", target.Scheme, target.Endpoint))

	return r, nil
}

func (r etcdResolver) Scheme() string {
	return schema
}

func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOption) {
}

func (r etcdResolver) Close() {
}

func (r *etcdResolver) watch(keyPrefix string) {
	var addrList []resolver.Address

	response, err := r.client.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err == nil {
		for i := range response.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(response.Kvs[i].Key), keyPrefix)})
		}
	}

	r.cc.NewAddress(addrList)

	rch := r.client.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.cc.NewAddress(addrList)
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.NewAddress(addrList)
				}
			}
		}
	}
}

func exist(address []resolver.Address, addr string) bool {
	for a := range address {
		if address[a].Addr == addr {
			return true
		}
	}
	return false
}

func remove(address []resolver.Address, addr string) ([]resolver.Address, bool) {
	for a := range address {
		if address[a].Addr == addr {
			address[a] = address[len(address)-1]
			return address[:len(address)-1], true
		}
	}
	return nil, false
}
