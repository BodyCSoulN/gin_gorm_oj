package service

import (
	"gin_gorm_oj/define"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("getProblemList Page strconv err : ", err)
		return
	}

	page = (page - 1) * size

	keyword := c.Query("keyword")

	data := make([]*models.ProblemBasic, 0)
	var count int64

	tx := models.GetProblemList(keyword)
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&data).Error

	if err != nil {
		log.Println("get problem list err : ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"list":  data,
		},
	})
}
