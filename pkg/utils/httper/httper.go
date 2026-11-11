package httper

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// Get sends a GET request and returns the response body as a string
func Get(url string, head map[string]string) (response string) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("failed to create GET request", zap.Error(err), zap.String("url", url))
		return ""
	}

	for k, v := range head {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("GET request failed", zap.Error(err), zap.String("url", url))
		return ""
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read response body", zap.Error(err))
		return ""
	}
	return string(result)
}

// PersonGet sends a GET request with a shorter timeout
func PersonGet(url string) (response string) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("failed to create PersonGet request", zap.Error(err), zap.String("url", url))
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("PersonGet request failed", zap.Error(err), zap.String("url", url))
		return ""
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read PersonGet response", zap.Error(err))
		return ""
	}
	return string(result)
}

// Post sends a POST request with the specified data and content type
func Post(url string, data []byte, contentType string, head map[string]string) (content string) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Error("failed to create POST request", zap.Error(err), zap.String("url", url))
		return ""
	}
	req.Header.Add("content-type", contentType)
	for k, v := range head {
		req.Header.Add(k, v)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("POST request failed", zap.Error(err), zap.String("url", url))
		return ""
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read POST response", zap.Error(err))
		return ""
	}
	return string(result)
}

// ZeroTierGet sends a GET request to ZeroTier API
func ZeroTierGet(url string, head map[string]string) (content string, code int) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("failed to create ZeroTier request", zap.Error(err), zap.String("url", url))
		return "", http.StatusInternalServerError
	}
	for k, v := range head {
		req.Header.Add(k, v)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("ZeroTier request failed", zap.Error(err), zap.String("url", url))
		return "", http.StatusInternalServerError
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read ZeroTier response", zap.Error(err))
		return "", code
	}
	return string(result), code
}

// 发送GET请求
// url:请求地址
// response:请求返回的内容
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
