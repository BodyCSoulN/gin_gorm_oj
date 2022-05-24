package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(36);" json:"identity"`                 // 唯一标识
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"` // 问题的唯一标识
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity;"`
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"` // 用户的唯一标识
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"`
	CodePath        string        `gorm:"column:code_path;type:varchar(255)" json:"code_path"` // 代码路径
	Status          int           `gorm:"column:status;type:tinyint(1);" json:"status"`
}

func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problem_identity, user_identity string, status int) *gorm.DB {
	tx := DB.Model(new(SubmitBasic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB { // 在preload中添加函数，对中间结果进行操作
		return db.Omit("content")
	}).Preload("UserBasic")
	if problem_identity != "" {
		tx.Where("problem_identity = ?", problem_identity)
	}

	if user_identity != "" {
		tx.Where("user_identity = ?", user_identity)
	}

	if status != -1 {
		tx.Where("status = ?", status)
	}

	return tx

}
