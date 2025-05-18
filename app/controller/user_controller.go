package controller

import (
	"gin_base/app/helper/db_helper"
	"gin_base/app/helper/exception_helper"
	"gin_base/app/helper/response_helper"
	"gin_base/app/model"
	"github.com/gin-gonic/gin"
)

// 获取登录用户信息
func GetUserInfo(c *gin.Context) {

	var user model.User

	uid, _ := c.Get("uid")
	err := db_helper.Db().Where("id=?", uid).Omit("password").First(&user)
	if err.Error != nil {
		exception_helper.CommonException("用户不存在")
	}

	response_helper.Success(c, "获取成功", user)
}
