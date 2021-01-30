package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"techstacks.cn/techstacks/dgraph"
	e "techstacks.cn/techstacks/error"
	"techstacks.cn/techstacks/utils"

	"github.com/gin-gonic/gin"
	"log"
)

type Auth struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

//DoAuth DoAuth
func DoAuth(c *gin.Context) {
	var auth Auth
	c.BindJSON(&auth)

	checkPassport := `
query checkPassword($email: String!, $password: String!) {
  checkPassportPassword(email: $email, password: $password) {
    email
    wxid
    author {
      domain
    }
  }
}
`
	//for error context
	eCtx := map[string]interface{}{
		"username": auth.Username,
	}

	resp, err := dgraph.DoRequire(checkPassport, &auth)
	if err != nil {
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_DGRAPH_FAILED,
			Error: ex,
		})
		log.Println("dgraph do require error:", ex.Error())
		return
	}
	defer resp.Body.Close()

	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_DGRAPH_FAILED,
			Error: ex,
		})
		log.Println("dgraph read user info error:", ex.Error())
		return
	}

	rspMap := map[string]interface{}{}
	err = json.Unmarshal(rspBody, &rspMap)
	if err != nil {
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_JSON_MARSHAL_FAILED,
			Error: ex,
		})
		log.Println("read user info error:", ex.Error())
		return
	}

	//检查有没有errors

	errors, _ := utils.GetMapPath(rspMap, "data", "checkPassportPassword", "errors")
	if errors != nil {
		b, exx := json.Marshal(errors)
		if exx != nil {
			b = []byte(exx.Error())
		}
		ex := e.NewError(string(b), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_DGRAPH_FAILED,
			Error: ex,
		})
		log.Println("Get usr info err:", ex.Error())
		return
	}

	//检查data中有没有checkPassportPassword
	checked, err := utils.GetMapPath(rspMap, "data", "checkPassportPassword")
	if err != nil {
		eCtx["rspBody"] = string(rspBody)
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_DGRAPH_FAILED,
			Error: ex,
		})
		log.Println("Get usr info err: ", ex.Error())
		return
	}

	if checked == nil {
		eCtx["rspBody"] = string(rspBody)
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_AUTH_FAILED,
			Error: ex,
		})
		log.Println("Check password failed!", ex)
		return
	}

	domain, err := utils.GetMapPath(rspMap, "data", "checkPassportPassword", "author", "domain")
	if err != nil {
		eCtx["rspMap"] = rspMap
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_USERDATA_ERROR,
			Error: ex,
		})
		log.Println("Get usr info err: ", ex.Error())
		return
	}
	domainS, ok := domain.(string)
	if !ok {
		ex := e.NewError("data  domain is not string", eCtx)
		c.JSON(http.StatusOK, Response{
			Code:  e.ERROR_USERDATA_ERROR,
			Error: ex,
		})
		log.Println("Get usr info err: ", ex.Error())
		return
	}
	token, err := utils.GenerateToken(auth.Username, domainS)
	if err != nil {
		ex := e.NewError(err.Error(), eCtx)
		c.JSON(http.StatusOK, &Response{
			Code:  e.ERROR_AUTH_TOKEN,
			Error: ex,
		})
		log.Println("Generate Token error:", ex.Error())
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
	rspDataM["token"] = token
	c.JSON(http.StatusOK, Response{
		Code: e.SUCCESS,
		Data: rspDataM,
	})
}
