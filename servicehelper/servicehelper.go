package servicehelper

import (
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"time"
)

// 服务启动
func ServiceStart(address, etcd_url, base_path string, server *server.Server) {
	// 注册etcd
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + address,
		EtcdServers:    []string{etcd_url},
		BasePath:       base_path,
		Metrics:        metrics.NewRegistry(),
		Services:       make([]string, 0),
		UpdateInterval: 30 * time.Second,
	}
	err := r.Start()
	if err != nil {
		log.Println(err)
	}
	server.Plugins.Add(r)
}
