package service

import (
	"gin_gorm_oj/define"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("getProblemList Page strconv err : ", err)
		return
	}

	page = (page - 1) * size // 起始位置
	var count int64
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	var submitList []models.SubmitBasic
	err = tx.Count(&count).Offset(page).Limit(size).Find(&submitList).Error

	if err != nil {
		log.Println("Get Submit List Error :", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Submit List Error :" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"list":  submitList,
		},
	})

}
