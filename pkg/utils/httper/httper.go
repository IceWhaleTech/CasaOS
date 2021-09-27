package httper

import (
	"bytes"
	"encoding/json"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

//发送GET请求
//url:请求地址
//response:请求返回的内容
func Get(url string, head map[string]string) (response string) {
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	req.BasicAuth()
	for k, v := range head {
		req.Header.Add(k, v)
	}
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		//需要错误日志的处理
		//loger.Error(error)
		return ""
		//panic(error)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			//loger.Error(err)
			return ""
			//	panic(err)
		}
	}
	response = result.String()
	return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func Post(url string, data interface{}, contentType string) (content string) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func ZeroTierPost(url string, data map[string]string, head map[string]string, cookies []*http.Cookie) (content string, code int) {
	b, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	for k, v := range head {
		req.Header.Add(k, v)
	}
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, error := client.Do(req)

	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	code = resp.StatusCode
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func ZeroTierPostJson(url string, data string, head map[string]string) (content string, code int) {
	var postData *bytes.Buffer

	jsonStr := []byte(data)
	postData = bytes.NewBuffer(jsonStr)

	req, err := http.NewRequest("POST", url, postData)
	for k, v := range head {
		req.Header.Add(k, v)
	}
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, error := client.Do(req)

	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	code = resp.StatusCode
	return
}

func ZeroTierDelete(url string, head map[string]string) (content string, code int) {

	req, err := http.NewRequest("DELETE", url, nil)
	for k, v := range head {
		req.Header.Add(k, v)
	}
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, error := client.Do(req)

	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	code = resp.StatusCode
	return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func ZeroTierGet(url string, head map[string]string) (content string, code int) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	for k, v := range head {
		req.Header.Add(k, v)
	}
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, error := client.Do(req)

	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	code = resp.StatusCode
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}

//发送GET请求
//url:请求地址
//response:请求返回的内容
func OasisGet(url string) (response string) {

	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	return Get(url, head)

}
