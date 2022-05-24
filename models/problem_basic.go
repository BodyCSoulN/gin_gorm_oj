package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(36);" json:"identity"`   // 问题表的唯一标识
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`                 // 关联问题分类表
	Title             string             `gorm:"column:title;type:varchar(255);" json:"title"`        // 题目
	Content           string             `gorm:"column:content;type:text;" json:"content"`            // 正文
	MaxRuntime        int                `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"` // 最大运行时间
	MaxMem            int                `gorm:"column:max_mem;type:int(11);" json:"max_mem"`         // 最大运行内存
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")

	// 根据categoryIdentity找到相应的问题
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (SELECT category_basic.id from category_basic WHERE category_basic.identity = ?)", categoryIdentity)
	}

	return tx
}
