package router

import (
	"github.com/gin-gonic/gin"
	"techstacks.cn/techstacks/api"
	"techstacks.cn/techstacks/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.POST("register", api.UserRegister)
	r.POST("login", api.DoAuth)
	r.POST("graphql", api.ReverseProxy("http://localhost:8080"))

	apiv1 := r.Group("/ping").Use(middleware.JWT())
	apiv1.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}
