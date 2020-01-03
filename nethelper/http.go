package nethelper

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

/**
 * 获取get请求数据
 */
func HttpGet(path string) string {
	resp, err := http.Get(path)
	if err != nil {
		log.Println("Http Get is error: ", err.Error())
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	return string(s)
}

/**
 * 获取post请求数据
 */
func HttpPost(path string, params map[string]string, header map[string]string) string {
	// 设置参数
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	// 建立连接
	body := ioutil.NopCloser(strings.NewReader(v.Encode())) // 把form 数据编下码
	client := &http.Client{}                                // 客户端
	request, err := http.NewRequest("POST", path, body)
	if err != nil {
		log.Println("Http Post is error: ", err.Error())
		return ""
	}
	// 设置头信息, Content-Type必须设定, POST参数才能正常提交
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	for k, val := range header {
		request.Header.Set(k, val)
	}

	resp, err := client.Do(request) // 发送请求
	defer func() {
		if err := recover(); err != nil {
			log.Printf("http.post error %v\n", err)
		}
	}()
	defer resp.Body.Close() // 一定要关闭resp.Body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Http Post is error: ", err.Error())
		return ""
	}
	return string(content)
}

/**
 * post请求，json请求数据
 */
func HttpPostJson(path string, param string, header map[string]string) string {
	jsonStr := []byte(param)
	body := bytes.NewBuffer(jsonStr) // 把form 数据编下码
	client := &http.Client{}
	request, err := http.NewRequest("POST", path, body)
	if err != nil {
		log.Println("http.post error ", err.Error())
		return ""
	}

	// 设置头信息, Content-Type必须设定, POST参数才能正常提交
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for k, val := range header {
		request.Header.Set(k, val)
	}

	resp, err := client.Do(request) // 发送请求
	defer func() {
		if err := recover(); err != nil {
			log.Printf("http.post error %v\n", err)
		}
	}()
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("http.post error ", err.Error())
		return ""
	}
	return string(content)
}

// 获取本地ip
func GetLocalhostIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
