package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"strings"
	"techstacks.cn/techstacks/api"
	"techstacks.cn/techstacks/dgraph"
	e "techstacks.cn/techstacks/error"
	"testing"
	"time"
)

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

//因为 dgraph是反向代理调用的， 会导致httptest.ResponseRecorder报错，封装一下,参考
//https://stackoverflow.com/questions/33968840/how-to-test-reverse-proxy-with-martini-in-go
type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

func TestRegisterAndLogin(t *testing.T) {
	auth := api.Auth{
		Username: RandString(5),
		Password: "1234qwer",
	}

	register := api.Register{
		Auth:   auth,
		Domain: auth.Username,
		Name:   auth.Username,
	}
	router := InitRouter()
	dgraph.Setup(&dgraph.Config{
		Hosts: []string{"http://localhost:8080/graphql"},
	})

	Convey("Test Register", t, func() {
		bs, err := json.Marshal(register)
		So(err, ShouldEqual, nil)
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(bs))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		result := w.Result()
		defer result.Body.Close()
		rspBody, err := ioutil.ReadAll(result.Body)
		So(err, ShouldEqual, nil)
		fmt.Println("rspBody = " + string(rspBody))

		var registRsp api.Response
		err = json.Unmarshal(rspBody, &registRsp)
		So(err, ShouldEqual, nil)
		So(registRsp.Code, ShouldEqual, e.SUCCESS)

	})

	Convey("Test Auth With JWT", t, func() {

		bs, err := json.Marshal(auth)
		So(err, ShouldEqual, nil)
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(bs))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		result := w.Result()
		defer result.Body.Close()

		rspBody, err := ioutil.ReadAll(result.Body)
		So(err, ShouldEqual, nil)
		var loginRsp api.Response
		err = json.Unmarshal(rspBody, &loginRsp)
		So(err, ShouldEqual, nil)
		So(loginRsp.Code, ShouldEqual, e.SUCCESS)

		rspData := loginRsp.Data
		//So(ok, ShouldEqual, true)
		token, ok := rspData["token"]
		So(ok, ShouldEqual, true)
		tokenS, ok := token.(string)
		So(ok, ShouldEqual, true)
		fmt.Println("tokenS = ", tokenS)

		Convey("Test Auth Api", func() {
			req := httptest.NewRequest("GET", "/ping", nil)
			req.Header["X-Auth-Token"] = []string{tokenS}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()

			rspBody, err := ioutil.ReadAll(result.Body)
			So(err, ShouldEqual, nil)
			fmt.Println("rspBody = " + string(rspBody))
			var pingRsp map[string]interface{}
			err = json.Unmarshal(rspBody, &pingRsp)
			So(err, ShouldEqual, nil)
			So(pingRsp["message"], ShouldEqual, "pong")
		})

		Convey("Test Add Blog", func() {
			q := `mutation {
	addBlog(input: [{ title: "abcde", text: "123456", author: { domain: "` + register.Domain + `" } }]) {
		blog {
			blogID
		}
	}
}`
			body := bytes.NewReader([]byte(q))
			req := httptest.NewRequest("POST", "/graphql", body)
			req.Header["X-Auth-Token"] = []string{tokenS}
			req.Header["Content-Type"] = []string{"application/graphql"}
			w := httptest.NewRecorder()
			//w := newCloseNotifyingRecorder() //反向代理特有的调用方式
			router.ServeHTTP(w, req)
			result := w.Result()
			defer result.Body.Close()

			rspBody, err := ioutil.ReadAll(result.Body)
			So(err, ShouldEqual, nil)
			fmt.Println("rspBody = " + string(rspBody))
		})
	})
}
