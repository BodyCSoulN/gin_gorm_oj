package service

import (
	"gin_gorm_oj/define"
	"gin_gorm_oj/helper"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user_identity"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "user identity is null",
		})
		return
	}

	userBaisc := new(models.UserBasic)
	err := models.DB.Omit("password").Where("identity = ?", identity).
		//Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&userBaisc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "current user not exists",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Get User Detail Error:" + err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": userBaisc,
	})

}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "username or password cant be null",
		})
		return
	}

	// md5 转换
	password = helper.GetMd5(password)

	data := new(models.UserBasic)

	err := models.DB.Where("name = ? and password = ?", username, password).First(&data).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "wrong username or password",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "User Login Error:" + err.Error(),
		})
		return
	}

	token, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate Token Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})

}

// SendVerifyCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /send-code [post]
func SendVerifyCode(c *gin.Context) {
	email := c.PostForm("email")

	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "email address is null",
		})
		return
	}

	verifyCode := helper.GenerateVerifyCode()
	err := models.RDB.Set(c, email, verifyCode, time.Minute*5).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Set Redis Key-Value Error:" + err.Error(),
		})
		return
	}

	err = helper.SendVerifyCode(email, verifyCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Send Email Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "successfully send an email to your mailbox. Please check it.",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param email formData string true "email"
// @Param verify-code formData string true "verify_code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	receiveCode := c.PostForm("verify-code")
	name := c.PostForm("name")
	password := c.PostForm("password")

	if email == "" || receiveCode == "" || name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "info requierd",
		})
		return
	}

	phone := c.PostForm("phone")

	// 检查邮箱是否已注册
	var count int64
	err := models.DB.Model(new(models.UserBasic)).Where("email = ?", email).Count(&count).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Info Error:" + err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "current email is exists",
		})
		return
	}

	// 验证码检查
	sendCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Verify Code Error:" + err.Error(),
		})
		return
	}
	if sendCode != receiveCode {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "verify code not match",
		})
		return
	}
	// 数据插入

	userIdentity := helper.GenerateUUid()
	data := &models.UserBasic{
		Identity: userIdentity,
		Name:     name,
		Password: helper.GetMd5(password),
		Phone:    phone,
		Email:    email,
	}

	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Create User Error:" + err.Error(),
		})
		return
	}

	token, err := helper.GenerateToken(userIdentity, name, data.IsAdmin)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate Token Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("getProblemList Page strconv err : ", err)
		return
	}
	page = (page - 1) * size // 起始位置
	var count int64
	data := make([]*models.UserBasic, 0)
	err = models.DB.Model(new(models.UserBasic)).Count(&count).
		Order("finish_problem_num DESC, submit_num ASC").
		Offset(page).Limit(size).Find(&data).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Get Rank List Error:" + err.Error(),
		})
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
