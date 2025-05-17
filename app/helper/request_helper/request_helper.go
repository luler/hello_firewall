package request_helper

import (
	"encoding/json"
	"gin_base/app/helper/helper"
	"gin_base/app/helper/valid_helper"
	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
	"net/url"
	"reflect"
	"strings"
)

// 获取请求参数-返回map类型
func Input(c *gin.Context, fields ...string) map[string]interface{} {
	// 获取所有的 GET 参数
	param1 := ParamGet(c, fields...)
	// 获取所有的 POST 参数
	param2 := ParamPostForm(c, fields...)
	// 获取multi-form类型所有的 POST 参数
	param3 := ParamMultipartForm(c, fields...)
	//获取json类型参数
	param4 := ParamRawJson(c, fields...)
	//合并参数
	param := helper.MergeMaps(param1, param2, param3, param4)
	//参数过滤
	param = helper.FilterMap(param, fields)
	return param
}

// 获取GET请求参数
func ParamGet(c *gin.Context, fields ...string) map[string]interface{} {
	return paramAll(c, "ParamGet", fields...)
}

// 获取PostForm请求参数
func ParamPostForm(c *gin.Context, fields ...string) map[string]interface{} {
	return paramAll(c, "ParamPostForm", fields...)
}

// 获取MultipartForm请求参数
func ParamMultipartForm(c *gin.Context, fields ...string) map[string]interface{} {
	return paramAll(c, "ParamMultipartForm", fields...)
}

// 获取json格式的请求参数
func ParamRawJson(c *gin.Context, fields ...string) map[string]interface{} {
	return paramAll(c, "ParamRawJson", fields...)
}

// 获取参数并验证
func InputStruct(c *gin.Context, param interface{}) {
	data := Input(c)
	mapstructure.Decode(data, param)
	valid_helper.Check(param)
}

// 获取参数并验证
func ParamGetStruct(c *gin.Context, param interface{}) {
	data := ParamGet(c)
	mapstructure.Decode(data, param)
	valid_helper.Check(param)
}

// 获取参数并验证
func ParamPostFormStruct(c *gin.Context, param interface{}) {
	data := ParamPostForm(c)
	mapstructure.Decode(data, param)
	valid_helper.Check(param)
}

// 获取参数并验证
func ParamMultipartFormStruct(c *gin.Context, param interface{}) {
	data := ParamMultipartForm(c)
	mapstructure.Decode(data, param)
	valid_helper.Check(param)
}

// 获取参数并验证
func ParamRawJsonStruct(c *gin.Context, param interface{}) {
	data := ParamRawJson(c)
	mapstructure.Decode(data, param)
	valid_helper.Check(param)
}

// 解析请求参数
func extractParam(c *gin.Context, queryParams url.Values, t string) map[string]interface{} {
	param := make(map[string]interface{})
	for key, values := range queryParams {
		if strings.HasSuffix(key, "]") {
			//存在类似param[1]的参数，需要特殊解析下
			key = key[0:strings.Index(key, "[")]
			switch t {
			case "get":
				qm := c.QueryMap(key)
				param[key] = qm
			case "post":
				qm := c.PostFormMap(key)
				param[key] = qm
			}

		} else {
			for _, value := range values {
				//判断是否存在多个相同参数，是就组成数组
				if _, ok := param[key]; ok {
					if reflect.TypeOf(param[key]).Kind() == reflect.Array {
						param[key] = append(param[key].([]interface{}), value)
					} else {
						param[key] = []interface{}{
							param[key],
							value,
						}
					}
				} else {
					param[key] = value
				}
			}
		}
	}

	return param
}

// 获取请求参数
func paramAll(c *gin.Context, method string, fields ...string) map[string]interface{} {
	data, exists := c.Get(method)
	param := make(map[string]interface{})
	if !exists {
		switch method {
		case "ParamGet":
			// 获取所有的 GET 参数
			queryParams := c.Request.URL.Query()
			for key, value := range extractParam(c, queryParams, "get") {
				param[key] = value
			}
		case "ParamPostForm":
			// 获取所有的PostForm参数
			c.Request.ParseForm()
			queryParams := c.Request.PostForm
			for key, value := range extractParam(c, queryParams, "post") {
				param[key] = value
			}
		case "ParamMultipartForm":
			// 获取multi-form类型所有的 POST 参数
			err := c.Request.ParseMultipartForm(32 << 24)
			if err == nil {
				queryParams := c.Request.MultipartForm.Value
				for key, value := range extractParam(c, queryParams, "post") {
					param[key] = value
				}
			}
		case "ParamRawJson":
			//获取json类型参数
			raw, _ := c.GetRawData()
			jsonData := make(map[string]interface{})
			json.Unmarshal(raw, &jsonData)
			for key, value := range jsonData {
				param[key] = value
			}
		}
		c.Set(method, param)
	} else {
		param = data.(map[string]interface{})
	}

	//参数过滤
	param = helper.FilterMap(param, fields)
	return param
}
