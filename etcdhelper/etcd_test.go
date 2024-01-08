package etcdhelper

import (
	"testing"
)

func TestEtcdGet(t *testing.T) {
	etcd := &EtcdConfig{
		Endpoints: []string{"192.168.2.192:2379"},
		UseCert:   false,
		CertFile:  "",
		KeyFile:   "",
		CaFile:    "",
		Timeout:   5,
		UseUser:   false,
		Username:  "",
		Password:  "",
	}
	//ret, err := etcd.EtcdGet("/oka/share/api")
	//fmt.Println(ret, err)

	//ret, err := etcd.EtcdGetWithPrefix("/oka/share")
	//fmt.Println(ret, err)

	//err := etcd.EtcdPut("/newKey", "xp xp ioioio")
	//fmt.Println(err)

	//err := etcd.EtcdPutWithLease("/newKey", "xp xp ioioio", 30, 0)
	//fmt.Println(err)

	etcd.etcdLock()
}
