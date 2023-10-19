package nethelper

import (
	"bytes"
	"errors"
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
	//body := ioutil.NopCloser(strings.NewReader(v.Encode())) // 把form 数据编下码
	body := strings.NewReader(v.Encode())
	client := &http.Client{} // 客户端
	request, err := http.NewRequest("POST", path, body)
	if err != nil {
		log.Println("Http Post is error: ", err.Error())
		return ""
	}
	// 设置头信息, Content-Type必须设定, POST参数才能正常提交
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		log.Println("http.NewRequest error ", err.Error())
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
		log.Println("http.ReadAll error ", err.Error())
		return ""
	}
	return string(content)
}

// 参数返回信息
type HttpResponse struct {
	Err    error  // 错误信息
	Body   string // 内容
	Status int    // 状态码
}

/**
 * http请求，综合
 * @param urlPath string 地址
 * @param param map[string]string 请求参数
 * @param options map[string]string 额外请求头部
 * @param isPost bool 是否是post
 * @return HttpResponse 结构数据
 */
func HttpRequest(urlPath string, param map[string]string, options map[string]string,
	isPost bool, reqBody string) HttpResponse {
	client := &http.Client{}
	resp := HttpResponse{
		Status: http.StatusInternalServerError,
	}
	err := errors.New("")
	var request *http.Request
	//提交请求
	if isPost {
		data := reqBody
		if reqBody == "" {
			requestInfo := url.Values{}
			for key, value := range param {
				requestInfo.Add(key, value)
			}
			data = requestInfo.Encode()
		}
		request, err = http.NewRequest("POST", urlPath, strings.NewReader(data))
	} else {
		request, err = http.NewRequest("GET", urlPath, nil)
	}
	if err != nil {
		resp.Err = err
		return resp
	}
	//设置请求头部
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"+
		" AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	for key, value := range options {
		request.Header.Set(key, value)
	}

	//处理返回结果
	response, err := client.Do(request)
	if err != nil {
		log.Println("error while request:" + urlPath + ":" + err.Error())
		resp.Body = ""
		resp.Status = 500
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error while request:" + urlPath + ":" + err.Error())
			resp.Body = ""
		} else {
			resp.Body = string(body)
			//返回的状态码
			resp.Status = response.StatusCode
		}
	}
	return resp
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
