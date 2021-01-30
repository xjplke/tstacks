package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"techstacks.cn/techstacks/dgraph"
	e "techstacks.cn/techstacks/error"
	"techstacks.cn/techstacks/utils"
)

type Register struct {
	Auth
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

//UserRegister Register
func UserRegister(c *gin.Context) {
	var register Register
	c.BindJSON(&register)

	addPassport := `
mutation addPassport($email: String!, $password: String!, $domain: String!, $name: String!) {
  addPassport(input: [{email: $email, password: $password, author:{domain: $domain, name: $name}}]) {
    passport {
      email
	  author {
      	domain
		name
      }
	}
  }
}
`
	//for error context
	eCtx := map[string]interface{}{
		"username": register.Username,
	}
	//TODO store Passport and Author
	resp, err := dgraph.DoRequire(addPassport, register)
	if err != nil {
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR,
			Error: ex,
		})
		log.Println("dgraph dorequire error:", ex.Error())
		return
	}
	defer resp.Body.Close()

	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil { //一般来说说这里不会报错，但是如果这里报错，实际注册是成功了，只是返回信息有点问题，应该如何处理？？？
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR,
			Error: ex,
		})
		log.Println("dgraph dorequire readbody error:", ex.Error())
		return
	}

	rspMap := map[string]interface{}{}

	err = json.Unmarshal(rspBody, &rspMap)
	if err != nil {
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR,
			Error: ex,
		})
		return
	}

	errors, _ := utils.GetMapPath(rspMap, "errors")
	if errors != nil {
		b, exx := json.Marshal(errors)
		if exx != nil {
			b = []byte(exx.Error())
		}
		ex := e.NewError(string(b), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR,
			Error: ex,
		})
		return
	}

	rspData, _ := utils.GetMapPath(rspMap, "data") //这里应该不会返回error，前面已经获取过了。

	rspDataM, ok := rspData.(map[string]interface{})
	if !ok {
		eCtx["rspData"] = rspData
		ex := e.NewError("Data Error", eCtx)
		c.JSON(http.StatusOK, &Response{
			Code:  e.ERROR_USERDATA_ERROR,
			Error: ex,
		})
		log.Println("Get usr info type err:", ex.Error())
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: e.SUCCESS,
		Data: rspDataM,
	})

	return
}
