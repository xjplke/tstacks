package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httputil"
	"net/url"
)

func ReverseProxy(target string) gin.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		log.Println("Reverse Proxy target url could not be parsed:", err)
		return nil
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func ProxyHandler(c *gin.Context) {

	c.Header()

}
