package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model
	Identity        string `gorm:"column:identity;type:varchar(36);" json:"identity"`                 // 唯一标识
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"` // 问题的唯一标识
	UserIdentity    string `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`       // 用户的唯一标识
	CodePath        string `gorm:"column:code_path;type:varchar(255)" json:"code_path"`               // 代码路径
}

func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}
