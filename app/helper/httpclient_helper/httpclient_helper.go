package log_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	defaultOptions http.Client
}

type HttpClientResponse struct {
	HttpCode     int
	Body         string
	ErrorMessage string
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		defaultOptions: http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Make a general request (GET/POST)
func (h *HttpClient) request(method, requestUrl string, params map[string]interface{}, headers map[string]string) *HttpClientResponse {
	client := h.defaultOptions
	httpClientResponse := &HttpClientResponse{
		HttpCode:     0,
		Body:         "",
		ErrorMessage: "",
	}
	// 准备请求数据
	var requestBody []byte
	if method == "POST" && len(params) > 0 {
		//判断是否json请求
		if headers["Content-Type"] == "application/json" {
			requestBody, _ = json.Marshal(params)
		} else {
			formData := url.Values{}
			for key, value := range params {
				formData.Add(key, fmt.Sprintf("%v", value))
			}
			requestBody = []byte(formData.Encode())
		}
	}

	// 构造请求
	req, err := http.NewRequest(method, requestUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		httpClientResponse.ErrorMessage = err.Error()
		return httpClientResponse
	}

	//设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 默认post请求内容类型
	if method == "POST" && len(params) > 0 && headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// get请求参数构造
	if method == "GET" && len(params) > 0 {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, fmt.Sprintf("%v", value))
		}
		req.URL.RawQuery = q.Encode()
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		httpClientResponse.ErrorMessage = err.Error()
		return httpClientResponse
	}
	defer resp.Body.Close()

	// 获取请求返回
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		httpClientResponse.ErrorMessage = err.Error()
		return httpClientResponse
	}

	httpClientResponse.HttpCode = resp.StatusCode
	httpClientResponse.Body = string(body)
	return httpClientResponse
}

// get请求
func (h *HttpClient) Get(requestUrl string, params map[string]interface{}, headers map[string]string) *HttpClientResponse {
	return h.request("GET", requestUrl, params, headers)
}

// post请求
func (h *HttpClient) Post(requestUrl string, params map[string]interface{}, headers map[string]string) *HttpClientResponse {
	return h.request("POST", requestUrl, params, headers)
}

// post+json请求
func (h *HttpClient) JsonPost(requestUrl string, data map[string]interface{}, headers map[string]string) *HttpClientResponse {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return h.request("POST", requestUrl, data, headers)
}
