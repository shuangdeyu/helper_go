package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

/**
 * etcd操作示例
 */

func TestEtcd() {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.2.192:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("client etcd error: " + err.Error())
	}
	defer etcdClient.Close()

	// 设置键值对
	key := "example-key"
	value := "example-value"
	_, err = etcdClient.Put(context.Background(), key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set key %s to value %s\n", key, value)

	// 获取键值对
	getResponse, err := etcdClient.Get(context.Background(), key)
	if err != nil {
		log.Fatal(err)
	}
	for _, kv := range getResponse.Kvs {
		fmt.Printf("Got key %s with value %s\n", kv.Key, kv.Value)
	}

	// 观察键的变化
	watchCh := etcdClient.Watch(context.Background(), key)
	for watchResp := range watchCh {
		for _, event := range watchResp.Events {
			fmt.Printf("Event Type: %v, Key: %s, Value: %s\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}
}
