package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	e "techstacks.cn/techstacks/error"
	"techstacks.cn/techstacks/utils"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := ""
		// https://jwt.io/introduction/
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			code = e.ERROR_NOT_AUTH
		} else {
			tokens := strings.Split(bearer, " ")
			if len(tokens) != 2 {
				code = e.INVALID_PARAMS
			} else {
				token = tokens[1]
			}
		}

		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			_, err := utils.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}