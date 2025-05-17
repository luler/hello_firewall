package model

import (
	"gin_base/app/helper/type_helper"
)

type User struct {
	Id        uint             `gorm:"primarykey;autoIncrement;comment:用户表" json:"id"`
	Name      string           `gorm:"type:varchar(50);not null;default:'';comment:账号;unique" json:"name"`
	Password  string           `gorm:"type:varchar(100);not null;default:'';comment:密码" json:"password"`
	Status    int8             `gorm:"not null;default:0;comment:状态，0-禁用，1-启用" json:"status"`
	CreatedAt type_helper.Time `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt type_helper.Time `gorm:"comment:更新时间" json:"updatedAt"`
}
