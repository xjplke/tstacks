package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"techstacks.cn/techstacks/dgraph"
	e "techstacks.cn/techstacks/error"
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
	//req := &http.Request{}
	//req.Header = c.Request.Header
	//req.Method = c.Request.Method
	//req.Body
	rsp, err := dgraph.DoProxy(c.Request)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR,
			Error: e.Error{},
		})
		log.Println("Do Proxy error!", err.Error())
		return
	}

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR,
			Error: e.Error{},
		})
		log.Println("Do Proxy ReadAll error!", err.Error())
		return
	}

	fmt.Println("Do Proxy response = ", string(rspBody))

	rspDataM := map[string]interface{}{}
	rspMap := map[string]interface{}{}
	err = json.Unmarshal(rspBody, &rspMap)
	c.JSON(http.StatusOK, Response{
		Code: e.SUCCESS,
		Data: rspDataM,
	})
}
