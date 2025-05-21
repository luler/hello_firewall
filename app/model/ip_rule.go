package model

import (
	"gin_base/app/helper/type_helper"
)

// IPRule 表示一条IP封禁规则
type IPRule struct {
	Id        uint             `gorm:"primarykey;autoIncrement;comment:iptable规则记录表" json:"id"`
	IP        string           `gorm:"type:varchar(50);not null;default:'';comment:ip" json:"ip"`
	Protocol  string           `gorm:"type:varchar(10);not null;default:'';comment:通信协议,tcp、udp、icmp" json:"protocol"`
	Port      int              `gorm:"type:int(10);not null;default:0;comment:封禁端口号" json:"port"`
	Reason    string           `gorm:"type:varchar(255);not null;default:'';comment:封禁原因" json:"reason"`
	Status    int8             `gorm:"type:tinyint(4);not null;default:0;comment:状态，0-禁用，1-启用" json:"status"`
	ExpiredAt int64            `gorm:"type:bigint;not null;default:0;comment:过期时间，0-无过期时间，>0-过期时间戳" json:"expired_at"`
	CreatedAt type_helper.Time `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt type_helper.Time `gorm:"comment:更新时间" json:"updatedAt"`
}
