package middleware

import "github.com/gin-gonic/gin"

func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: check if user is admin

	}
}
