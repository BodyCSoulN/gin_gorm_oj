package service

import (
	"encoding/json"
	"gin_gorm_oj/define"
	"gin_gorm_oj/helper"
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

// CreateProblem
// @Tags 私有方法
// @Summary 创建问题
// @Param token header string true "token"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_memory formData int false "max_memory"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array true "test_cases"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /create-problem [post]
func CreateProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	max_runtime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	max_memory, _ := strconv.Atoi(c.PostForm("max_memory"))
	category_ids := c.PostFormArray("category_ids")
	test_cases := c.PostFormArray("test_cases")

	if title == "" || content == "" ||
		len(category_ids) == 0 || len(test_cases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "params can not be null",
		})

		return
	}

	identity := helper.GenerateUUid()

	data := models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: max_runtime,
		MaxMem:     max_memory,
	}

	// 处理分类
	ProblemCategories := make([]*models.ProblemCategory, 0)
	for _, id := range category_ids {
		categoryId, _ := strconv.Atoi(id)
		ProblemCategories = append(ProblemCategories, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(categoryId),
		})
	}
	data.ProblemCategories = ProblemCategories

	// 处理测试用例
	testCases := make([]*models.TestCase, 0)
	for _, testCase := range test_cases {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Wrong Format of Test Case" + err.Error(),
			})
			return
		}

		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Wrong Format of Test Case" + err.Error(),
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Wrong Format of Test Case" + err.Error(),
			})
			return
		}

		testCases = append(testCases, &models.TestCase{
			Identity:        helper.GenerateUUid(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		})
	}

	data.TestCases = testCases

	// 创建问题
	err := models.DB.Create(data).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Create Problem Error:" + err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{},
	})
}
