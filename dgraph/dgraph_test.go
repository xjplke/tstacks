package dgraph

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"techstacks.cn/techstacks/utils"
	"testing"
	"time"
)

type Auth struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}

func TestDgraphHttp(t *testing.T) {
	Setup(&Config{
		Hosts: []string{"http://localhost:8080/graphql"},
	})
	param := Auth{
		Username: RandString(5),
		Password: "1234qwer",
	}

	Convey("Test Register", t, func() {

		addPassport := `
mutation addPassport($email: String!, $password: String!) {
  addPassport(input: [{email: $email, password: $password, author:{domain: $email, name:$email}}]) {
    passport {
      email
	}
  }
}
`

		resp, err := DoRequire(addPassport, param)
		So(err, ShouldEqual, nil)
		defer resp.Body.Close()
		rspBody, err := ioutil.ReadAll(resp.Body)
		So(err, ShouldEqual, nil)

		rspMap := map[string]interface{}{}
		err = json.Unmarshal(rspBody, &rspMap)
		So(err, ShouldEqual, nil)
		errors, err := utils.GetMapPath(rspMap, "errors")
		So(err, ShouldNotEqual, nil)
		So(errors, ShouldEqual, nil)

		fmt.Println(string(rspBody))
	})

	Convey("Test login", t, func() {
		/*
		   		checkPassport := `
		   query do($email: String!, $password: String!) {
		     checkPassportPassword(input: [{email: $email, password: $password}]) {
		       email
		       wxid
		       author {
		         id
		       }
		     }
		   }
		   `
		*/

		checkPassport := `query {
  checkPassportPassword(email: "` + param.Username + `", password: "` + param.Password + `") {
    email
	author {
		domain
		name
	}
  }
}
`

		resp, err := DoRequire(checkPassport, nil)
		So(err, ShouldEqual, nil)
		defer resp.Body.Close()
		rspBody, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(rspBody))
		rspMap := map[string]interface{}{}
		err = json.Unmarshal(rspBody, &rspMap)
		So(err, ShouldEqual, nil)
		_, err = utils.GetMapPath(rspMap, "errors")
		So(err, ShouldNotEqual, nil)
		domain, err := utils.GetMapPath(rspMap, "data", "checkPassportPassword", "author", "domain")
		So(domain, ShouldEqual, param.Username)
	})
}
