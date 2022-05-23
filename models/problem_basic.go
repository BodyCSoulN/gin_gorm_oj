package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`   // 问题表的唯一标识
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`        // 题目
	Content    string `gorm:"column:content;type:text;" json:"content"`            // 正文
	MaxRuntime int    `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"` // 最大运行时间
	MaxMem     int    `gorm:"column:max_mem;type:int(11);" json:"max_mem"`         // 最大运行内存
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword string) *gorm.DB {
	return DB.Model(new(ProblemBasic)).
		Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")
	//data := make([]Problem, 10)
	//DB.Find(&data)
	//
	//for _, v := range data {
	//	fmt.Printf("problem ==> %v", v)
	//}
}
