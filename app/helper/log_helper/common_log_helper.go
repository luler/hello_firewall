package log_helper

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gin_base/app/helper/cache_helper"
	"github.com/syyongx/php2go"
	"io"
	"net/http"
	"os"
	"time"
)

// 推送日志数据到队列
func QueueCommonLog(logData []map[string]interface{}) {
	if os.Getenv("COMMON_LOG_ENABLE") != "true" || len(logData) == 0 {
		return
	}
	logData = defaultLogData(logData)

	for _, value := range logData {
		json_data, _ := json.Marshal(value)
		cache_helper.RedisHelper().Client.LPush(context.Background(), "CommonLog:Gin", string(json_data))
	}
}

// 数据格式处理
func defaultLogData(logData []map[string]interface{}) []map[string]interface{} {
	for _, value := range logData {
		if php2go.Empty(value["level"]) {
			value["level"] = "info"
		}
		if php2go.Empty(value["code"]) {
			value["code"] = 0
		}
		if php2go.Empty(value["url"]) {
			value["url"] = "log"
		}
		if php2go.Empty(value["waste_time"]) {
			value["waste_time"] = 0
		}
		if php2go.Empty(value["message"]) {
			value["message"] = ""
		}
		if php2go.Empty(value["other"]) {
			value["other"] = ""
		}
		if php2go.Empty(value["other_id"]) {
			value["other_id"] = ""
		}
		if php2go.Empty(value["create_time"]) {
			value["create_time"] = time.Now().UnixMilli()
		}
		if php2go.Empty(value["client_ip"]) {
			value["client_ip"] = "0.0.0.0"
		}
		if php2go.Empty(value["server_ip"]) {
			value["server_ip"] = "0.0.0.0"
		}
		value["project_name"] = os.Getenv("COMMON_LOG_PROJECT_NAME")
	}

	return logData
}

// 写日志
func SaveCommonLog(logData []map[string]interface{}) {
	if os.Getenv("COMMON_LOG_ENABLE") != "true" || len(logData) == 0 {
		return
	}

	logData = defaultLogData(logData)

	data := map[string]interface{}{
		"authorization": getCommonLogAccessToken(),
		"data":          logData,
	}
	url := os.Getenv("COMMON_LOG_HOST") + "/api/saveLog"

	jsonData, _ := json.Marshal(data)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	resp, _ := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if int(result["code"].(float64)) != 200 {
		fmt.Println("日志保存异常", result)
	}
}

// 获取授权码
func getCommonLogAccessToken() string {
	var accessToken string
	cacheKey := "GetCommonLogAccessToken"
	accessToken, _ = cache_helper.RedisHelper().RedisGet(cacheKey)
	if accessToken == "" {
		url := os.Getenv("COMMON_LOG_HOST") + "/api/getAccessToken"
		param := map[string]string{
			"appid":     os.Getenv("COMMON_LOG_APPID"),
			"appsecret": os.Getenv("COMMON_LOG_APPSECRET"),
		}
		jsonData, _ := json.Marshal(param)
		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}
		resp, _ := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		result = result["info"].(map[string]interface{})
		accessToken = result["access_token"].(string)

		cache_helper.RedisHelper().RedisSet(cacheKey, accessToken, time.Duration(int(result["expires_in"].(float64))-10)*time.Second)
	}

	return accessToken
}

// 推送日志
func PushCommonLog() {
	id := cache_helper.RedisHelper().RedisLock("PushCommonLog", time.Second*30)
	if id == "" {
		return
	}
	defer cache_helper.RedisHelper().RedisUnLock("PushCommonLog", id)

	var results []map[string]interface{}
	//一次最多处理500条
	for i := 0; i < 500; i++ {
		data, err := cache_helper.RedisHelper().Client.RPop(context.Background(), "CommonLog:Gin").Result()
		if err != nil {
			break
		}
		var result map[string]interface{}
		json.Unmarshal([]byte(data), &result)
		results = append(results, result)
	}

	SaveCommonLog(results)
}
