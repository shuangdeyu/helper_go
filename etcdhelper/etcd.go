package etcdhelper

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"time"
)

/**
 * etcd:是一个分布式键值对存储，设计用来可靠而快速的保存关键数据并提供访问。通过分布式锁，leader选举和写屏障(write barriers)来实现
 * 可靠的分布式协作。etcd集群是为高可用，持久性数据存储和检索而准备
 *
 * 功能概览：
 * 1. 存储 key-value 数据
 * 2. 监听数据变化
 * 3. 租约，续期
 * 4. 网关
 * 功能应用：
 * 1. KV数据库
 * 2. 服务注册发现
 * 3. 共享配置
 * 4. 协调分布式
 * 5. 分布式锁
 */

type EtcdConfig struct {
	Endpoints []string `json:"endpoints"`
	Timeout   int      `json:"timeout"`
	UseUser   bool     `json:"useUser"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	UseCert   bool     `json:"useCert"`
	CertFile  string   `json:"certFile"` // 用于与etcd的SSL/TLS连接的证书
	KeyFile   string   `json:"keyFile"`  // 证书的密钥。必须未加密
	CaFile    string   `json:"caFile"`   // 受信任的证书颁发机构
}

var EtcdClient *clientv3.Client

// InitEtcdClient 初始化etcd客户端
func (e *EtcdConfig) InitEtcdClient() (*clientv3.Client, error) {
	if EtcdClient == nil {
		// 创建配置项
		config := clientv3.Config{
			Endpoints: e.Endpoints,
		}
		if e.Timeout <= 0 {
			e.Timeout = 5
		}
		config.DialTimeout = time.Duration(e.Timeout) * time.Second
		if e.UseUser {
			config.Username = e.Username
			config.Password = e.Password
		}
		if e.UseCert {
			tlsInfo := transport.TLSInfo{
				CertFile:      e.CertFile,
				KeyFile:       e.KeyFile,
				TrustedCAFile: e.CaFile,
			}
			_tlsConfig, err := tlsInfo.ClientConfig()
			if err != nil {
				fmt.Printf("etccd tlsconfig failed, err:%v\n", err)
				return EtcdClient, err
			}
			config.TLS = _tlsConfig
		}

		// 创建客户端连接
		client, err := clientv3.New(config)
		if err != nil {
			fmt.Printf("etcd client failed, err:%v\n", err)
			return EtcdClient, err
		}
		EtcdClient = client
	}
	return EtcdClient, nil
}

// EtcdGet 获取键值最新版本的数据
func (e *EtcdConfig) EtcdGet(key string) (string, error) {
	cli, err := e.InitEtcdClient()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		fmt.Printf("get etcd failed, err:%v\n", err)
		return "", err
	}
	for k, ev := range resp.Kvs {
		if k == 0 {
			return string(ev.Value), nil
		}
	}
	return "", errors.New("key not found")
}

// EtcdGetWithPrefix 根据key前缀获取数据
func (e *EtcdConfig) EtcdGetWithPrefix(key string) (map[string]string, error) {
	ret := make(map[string]string)
	cli, err := e.InitEtcdClient()
	if err != nil {
		return ret, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		fmt.Printf("get etcd failed, err:%v\n", err)
		return ret, err
	}
	for _, ev := range resp.Kvs {
		ret[string(ev.Key)] = string(ev.Value)
	}
	return ret, nil
}

// EtcdPut 添加key-value
func (e *EtcdConfig) EtcdPut(key, value string) error {
	cli, err := e.InitEtcdClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return err
	}
	return nil
}

// EtcdPutWithLease 带租约添加key-value
func (e *EtcdConfig) EtcdPutWithLease(key, value string, expireTime, leaseId int64) error {
	cli, err := e.InitEtcdClient()
	if err != nil {
		return err
	}

	// 创建租约
	id := clientv3.LeaseID(leaseId)
	if leaseId == 0 && expireTime > 0 {
		resp, err := cli.Grant(context.TODO(), expireTime)
		if err != nil {
			return err
		}
		id = resp.ID
	}

	// 添加键值
	_, err = cli.Put(context.TODO(), key, value, clientv3.WithLease(id))
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return err
	}
	return nil
}

// watch监听操作，放在业务流程中进行监听，配合业务逻辑使用
// 这里放一个demo
func (e *EtcdConfig) etcdWatch(key string) {
	cli, err := e.InitEtcdClient()
	if err != nil {
		return
	}

	rch := cli.Watch(context.Background(), key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

// KeepAlive 也是一个demo，请配合业务逻辑使用
func (e *EtcdConfig) etcdKeepAlive(key string) {
	cli, err := e.InitEtcdClient()
	if err != nil {
		return
	}

	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}
	_, err = cli.Put(context.TODO(), "/lmh/", "lmh", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	for {
		ka := <-ch
		fmt.Println("ttl:", ka.TTL)
	}
}

// 分布式锁demo
func (e *EtcdConfig) etcdLock() {
	cli, err := e.InitEtcdClient()
	if err != nil {
		return
	}

	// 创建两个单独的会话用来演示锁竞争
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, "/my-lock/")

	s2, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s2.Close()
	m2 := concurrency.NewMutex(s2, "/my-lock/")

	// 会话s1获取锁
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// 等待直到会话s1释放了/my-lock/的锁
		if err := m2.Lock(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for s1")

	<-m2Locked
	fmt.Println("acquired lock for s2")
}
