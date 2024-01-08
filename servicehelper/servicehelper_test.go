package servicehelper

import (
	"context"
	"fmt"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"testing"
	"time"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func TestServiceStart(t *testing.T) {
	// etcd 3 的注册
	s := server.NewServer()

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@127.0.0.1:8972",
		EtcdServers:    []string{"192.168.2.192:2379"},
		BasePath:       "/rpcx_test",
		Metrics:        metrics.NewRegistry(),
		Services:       make([]string, 0),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		return
	}
	s.Plugins.Add(r)

	s.RegisterName("Arith", new(Arith), "")
	go s.Serve("tcp", "127.0.0.1:8972")
	defer s.Close()

	if len(r.Services) != 1 {
		t.Fatal("failed to register services in etcd")
	}
	fmt.Println(r.Services)

	if err := r.Stop(); err != nil {
		t.Fatal(err)
	}
}
