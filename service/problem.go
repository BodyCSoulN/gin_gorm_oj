package service

import (
	"gin_gorm_oj/define"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("getProblemList Page strconv err : ", err)
		return
	}

	page = (page - 1) * size // 起始位置

	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")

	data := make([]*models.ProblemBasic, 0)
	var count int64

	tx := models.GetProblemList(keyword, categoryIdentity)
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem_identity"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "problem identity is null",
		})

		return
	}

	problemBaisc := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&problemBaisc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "current problem not exists",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Get Problem Detail Error:" + err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": problemBaisc,
	})

}
