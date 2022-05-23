package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"` // 用户的唯一标识
	Name     string `gorm:"column:name;type:varchar(100)" json:"name"`         // 用户名
	Password string `gorm:"column:password;type:varchar(32)" json:"password"`  // 用户密码
	Phone    string `gorm:"column:phone;type:varchar(20)" json:"phone"`        // 手机号
	Email    string `gorm:"column:email;type:varchar(100)" json:"email"`       // 邮箱
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
