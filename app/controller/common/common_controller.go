package common

import (
	"gin_base/app/helper/response_helper"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	response_helper.Success(c, "访问成功")
}
