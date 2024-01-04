package servicehelper

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"

	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/codec"
)

const (
	XVersion           = "X-RPCX-Version"
	XMessageType       = "X-RPCX-MesssageType"
	XHeartbeat         = "X-RPCX-Heartbeat"
	XOneway            = "X-RPCX-Oneway"
	XMessageStatusType = "X-RPCX-MessageStatusType"
	XSerializeType     = "X-RPCX-SerializeType"
	XMessageID         = "X-RPCX-MessageID"
	XServicePath       = "X-RPCX-ServicePath"
	XServiceMethod     = "X-RPCX-ServiceMethod"
	XMeta              = "X-RPCX-Meta"
	XErrorMessage      = "X-RPCX-ErrorMessage"
)

type Out struct {
	OutMsg  interface{} `json:"outMsg"`
	OutData interface{} `json:"outData"`
}
type Output struct {
	Out Out
}

/**
 * 通过网关调用服务
 */
func CallService(addr, service, method string, params map[string]interface{}) (Out, error) {
	cc := &codec.MsgpackCodec{}
	data, _ := cc.Encode(params)

	req, err := http.NewRequest("POST", "http://"+addr, bytes.NewReader(data))
	if err != nil {
		log.Println(service+"."+method, " failed to create request: ", err)
		return Out{}, err
	}

	h := req.Header
	h.Set(XMessageID, "10000")
	h.Set(XMessageType, "0")
	h.Set(XSerializeType, "3")
	h.Set(XServicePath, service)
	h.Set(XServiceMethod, method)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(service+"."+method, " failed to call: ", err)
		return Out{}, err
	}
	defer res.Body.Close()

	replyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(service+"."+method, " failed to read response: ", err)
		return Out{}, err
	}

	output := &Output{}
	err = cc.Decode(replyData, output)
	if err != nil {
		log.Println(service+"."+method, " failed to decode reply: ", err, replyData)
		return Out{}, err
	}

	return output.Out, nil
}

/**
 * 通过etcd调用服务
 */
func CallServiceByEtcd(addr, path, service, method string, params map[string]interface{}) (Out, error) {
	d, _ := etcd_client.NewEtcdDiscovery(path, service, []string{addr}, false, nil)
	xclient := client.NewXClient(service, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	output := &Output{}
	err := xclient.Call(context.Background(), method, params, output)
	if err != nil {
		log.Println("failed to call CallServiceByEtcd: ", err)
		return Out{}, err
	}

	return output.Out, nil
}

/**
 * 通过TCP调用服务
 */
func CallServiceByTcp(addr, service, method string, params map[string]interface{}) (Out, error) {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	if err != nil {
		log.Println("failed to connection CallServiceByTcp: ", err)
		return Out{}, err
	}
	xclient := client.NewXClient(service, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	output := &Output{}
	err = xclient.Call(context.Background(), method, params, output)
	if err != nil {
		log.Println("failed to call CallServiceByTcp: ", err)
		return Out{}, err
	}

	return output.Out, nil
}
