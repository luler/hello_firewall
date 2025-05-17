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

// 封禁ip接口
func BanIp(c *gin.Context) {
	type Param struct {
		Ips      []string `validate:"required,min=1,max=100" label:"ip数组"`
		Protocol string   `validate:"required,oneof=tcp udp icmp" label:"通信协议"`
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
		Ips      []string `validate:"required,min=1,max=100" label:"ip数组"`
		Protocol string   `validate:"omitempty,oneof=tcp udp icmp" label:"通信协议"`
		Port     int      `validate:"omitempty,gte=1,lte=65535" label:"封禁端口号"`
	}
	var param Param
	request_helper.ParamRawJsonStruct(c, &param)

	// 动态构建查询条件
	query := db_helper.Db().Where("ip IN ?", param.Ips)

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
	//删除后查询所有规则，重置iptables规则
	var rules []*model.IPRule
	db_helper.Db().Find(&rules)

	err = logic.GetIPTablesManager().RebuildRules(rules)

	if err != nil {
		exception_helper.CommonException(err.Error())
	}
	response_helper.Success(c, "解封IP成功")
}

// 解封ip接口
func GetBanIpList(c *gin.Context) {
	type Param struct {
		Ip string `validate:"omitempty" label:"ip关键字"`
	}
	var param Param
	request_helper.ParamRawJsonStruct(c, &param)

	// 动态构建查询条件
	query := db_helper.Db().Model(&model.IPRule{})

	if param.Ip != "" {
		query = query.Where("ip LIKE ?", "%"+param.Ip+"%")
	}

	query = query.Order("id desc")
	//执行数据库删除
	data := db_helper.AutoPage(c, query)
	response_helper.Success(c, "获取成功", data)
}
