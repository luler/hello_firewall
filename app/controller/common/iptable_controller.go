package common

import (
	"fmt"
	"gin_base/app/helper/db_helper"
	"gin_base/app/helper/exception_helper"
	"gin_base/app/helper/request_helper"
	"gin_base/app/helper/response_helper"
	"gin_base/app/logic"
	"gin_base/app/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 封禁ip接口
// @Description  封禁ip接口
// @Tags IP封禁相关接口
// @Accept x-www-form-urlencoded
// @Produce  json
// @Param ips formData array true "ip数组，格式:[127.0.0.1]"
// @Param protocol formData string false "封禁协议,tcp udp icmp"
// @Param port formData int false "封禁端口号"
// @Param reason formData string false "封禁原因"
// @Success 200
// @Router /api/banIp [post]
func BanIp(c *gin.Context) {
	type Param struct {
		Ips      []string `validate:"required,min=1,max=100" label:"ip数组"`
		Protocol string   `validate:"omitempty,oneof=tcp udp icmp" label:"封禁协议"`
		Port     int      `validate:"omitempty,gte=1,lte=65535" label:"封禁端口号"`
		Minute   int      `validate:"omitempty,gte=0" label:"封禁时长分钟"`
		Reason   string   `validate:"omitempty,max=255" label:"封禁原因"`
	}
	var param Param
	request_helper.ParamRawJsonStruct(c, &param)
	if param.Protocol == "icmp" { //icmp不支持设置端口
		param.Port = 0
	}

	err := db_helper.Db().Transaction(func(tx *gorm.DB) error {
		for _, ip := range param.Ips {
			//存在则激活规则
			var existingRule model.IPRule
			err := tx.Where("ip = ? AND protocol = ? AND port = ?", ip, param.Protocol, param.Port).First(&existingRule).Error
			if err == nil {
				// Rule exists, update status to 1 if it was 0
				if existingRule.Status == 0 {
					existingRule.Status = 1
					existingRule.Reason = param.Reason
					if err := tx.Save(&existingRule).Error; err != nil {
						return fmt.Errorf("更新IP规则状态失败: %v", err)
					}
					if err := logic.GetIPTablesManager().ApplyRule(&existingRule); err != nil {
						return fmt.Errorf("设置iptables规则失败: %v", err)
					}
				}
				continue //处理完了，继续下一个
			}
			if err != gorm.ErrRecordNotFound {
				return fmt.Errorf("查询IP规则失败: %v", err)
			}
			//不存在就新增规则
			ipRule := model.IPRule{
				IP:       ip,
				Protocol: param.Protocol,
				Port:     param.Port,
				Status:   1,
				Reason:   param.Reason,
			}
			if err := tx.Save(&ipRule).Error; err != nil {
				return fmt.Errorf("保存IP规则失败: %v", err)
			}
			if err := logic.GetIPTablesManager().ApplyRule(&ipRule); err != nil {
				return fmt.Errorf("设置iptables规则失败: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		exception_helper.CommonException(err.Error())
	}
	response_helper.Success(c, "封禁IP成功")
}

// 解封ip接口
func UnBanIp(c *gin.Context) {
	type Param struct {
		Ids      []int    `validate:"omitempty,min=1,max=100" label:"id数组"`
		Ips      []string `validate:"omitempty,min=1,max=100" label:"ip数组"`
		Protocol string   `validate:"omitempty,oneof=tcp udp icmp" label:"通信协议"`
		Port     int      `validate:"omitempty,gte=1,lte=65535" label:"封禁端口号"`
	}
	var param Param
	request_helper.ParamRawJsonStruct(c, &param)
	if len(param.Ids) == 0 && len(param.Ips) == 0 {
		exception_helper.CommonException("id数组和ip数组不能同时为空")
	}
	// 动态构建查询条件
	query := db_helper.Db()
	if len(param.Ids) > 0 {
		query = query.Where("id IN ?", param.Ids)
	} else {
		query = query.Where("ip IN ?", param.Ips)
	}

	if param.Protocol != "" {
		query = query.Where("protocol = ?", param.Protocol)
	}

	if param.Port > 0 {
		query = query.Where("port = ?", param.Port)
	}
	//执行数据库删除
	err := query.Delete(&model.IPRule{}).Error
	if err != nil {
		exception_helper.CommonException(fmt.Sprintf("删除IP规则失败: %v", err))
	}

	err = logic.GetIPTablesManager().RebuildRules()
	if err != nil {
		exception_helper.CommonException(err.Error())
	}
	response_helper.Success(c, "解封IP成功")
}

// 解封ip接口
func ChangeStatus(c *gin.Context) {
	type Param struct {
		Id     int `validate:"required" label:"id"`
		Status int `validate:"omitempty,oneof=0 1" label:"状态"`
	}
	var param Param
	request_helper.ParamRawJsonStruct(c, &param)

	db_helper.Db().Transaction(func(tx *gorm.DB) error {
		// 动态构建查询条件
		var ipRule model.IPRule
		if err := tx.Where("id = ?", param.Id).First(&ipRule).Error; err != nil {
			exception_helper.CommonException(fmt.Sprintf("查询IP规则失败: %v", err))
		}
		ipRule.Status = int8(param.Status)
		//执行数据库更新
		if err := tx.Save(&ipRule).Error; err != nil {
			exception_helper.CommonException(fmt.Sprintf("更新IP规则失败: %v", err))
		}

		err := logic.GetIPTablesManager().RebuildRules()
		if err != nil {
			exception_helper.CommonException(err.Error())
		}
		return nil
	})

	response_helper.Success(c, "解封IP成功")
}

// 解封ip接口
func GetBanIpList(c *gin.Context) {
	type Param struct {
		Ip     string `validate:"omitempty" label:"ip关键字"`
		Status string `validate:"omitempty" label:"状态"`
	}
	var param Param
	request_helper.ParamGetStruct(c, &param)

	// 动态构建查询条件
	query := db_helper.Db().Model(&model.IPRule{})

	if param.Ip != "" {
		query = query.Where("ip LIKE ?", "%"+param.Ip+"%")
	}
	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	query = query.Order("id desc")
	//执行数据库删除
	data := db_helper.AutoPage(c, query)
	response_helper.Success(c, "获取成功", data)
}
