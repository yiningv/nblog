package router

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/pub/util"
	"net/http"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := OK
		auth := c.GetHeader("auth")
		token := strings.Split(auth, " ")

		if auth == "" {
			code = ERR
		} else {
			_, err := util.ParseToken(token[1])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = ERR
				default:
					code = ERR
				}
			}
		}

		if code != OK {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		c.Next()
	}
}
