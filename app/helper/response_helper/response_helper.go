package response_helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, message string, data ...interface{}) {
	res := make(map[string]interface{})
	res["code"] = http.StatusOK
	res["message"] = message
	res["data"] = []int{}
	if len(data) > 0 {
		res["data"] = data[0]
	}
	jsonResponse(c, res)
}

func Fail(c *gin.Context, message string, data ...interface{}) {
	res := make(map[string]interface{})
	res["code"] = http.StatusBadRequest
	res["message"] = message
	res["data"] = []int{}
	if len(data) > 0 {
		res["data"] = data[0]
	}
	jsonResponse(c, res)
}

func Common(c *gin.Context, code int, message string, data ...interface{}) {
	res := make(map[string]interface{})
	res["code"] = code
	res["message"] = message
	res["data"] = []int{}
	if len(data) > 0 {
		res["data"] = data[0]
	}
	jsonResponse(c, res)
}

func jsonResponse(c *gin.Context, res map[string]interface{}) {
	//保存参数
	c.Set("response_data", res)
	c.JSON(res["code"].(int), res)
}
