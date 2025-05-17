package db_helper

import (
	"fmt"
	"gin_base/app/helper/request_helper"
	"gin_base/app/helper/type_helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// 自动分页获取
func AutoPage(c *gin.Context, db *gorm.DB, isGetTotal_Page_PageSize ...int) map[string]interface{} {
	isGetTotal := 1
	Page := 0
	PageSize := 0
	if len(isGetTotal_Page_PageSize) >= 1 {
		isGetTotal = isGetTotal_Page_PageSize[0]
	}
	if len(isGetTotal_Page_PageSize) >= 2 {
		Page = isGetTotal_Page_PageSize[1]
	}
	if len(isGetTotal_Page_PageSize) >= 3 {
		PageSize = isGetTotal_Page_PageSize[2]
	}
	if Page == 0 {
		type Param struct {
			Page      any
			Page_Size any
		}
		var param Param
		request_helper.InputStruct(c, &param)

		if p, err := strconv.Atoi(fmt.Sprintf("%v", param.Page)); err == nil {
			Page = p
		} else {
			Page = 1
		}
		switch param.Page_Size.(type) {
		case string:
			PageSize, _ = strconv.Atoi(param.Page_Size.(string))
		case int, int8, int16, int32, int64:
			PageSize = param.Page_Size.(int)
		default:
			PageSize = 10
		}
	}

	var total int64
	if isGetTotal == 1 {
		db.Count(&total)
	}
	var data []map[string]interface{}
	if isGetTotal != 1 || (isGetTotal == 1 && total > 0) {
		if Page == -1 || PageSize == -1 { //不分页
			db.Find(&data)
		} else { //分页
			offset := (Page - 1) * PageSize
			db.Limit(PageSize).Offset(offset).Find(&data)
		}
	}
	if data == nil {
		data = []map[string]interface{}{}
	}
	//字段处理
	for _, item := range data {
		for key, value := range item {
			// 将字段名转换为大驼峰
			if tt, ok := value.(time.Time); ok {
				value = type_helper.Time(tt)
			}
			item[key] = value
		}
	}
	res := map[string]interface{}{}
	res["list"] = data
	if Page != -1 && PageSize != -1 { //分页
		res["page"] = Page
		res["page_size"] = PageSize
	}
	if isGetTotal == 1 {
		res["total"] = total
	}

	return res
}
