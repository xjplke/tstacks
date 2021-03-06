package dgraph

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Config struct {
	Hosts []string
}

type Client struct {
	hosts      []string //集群多个节点的时候也可以
	httpClient http.Client
}

var client *Client

func Setup(c *Config) {
	client = &Client{
		hosts:      c.Hosts,
		httpClient: http.Client{},
	}
}

func (c *Client) DoRequire(q string, param interface{}) (*http.Response, error) {
	qx := map[string]interface{}{
		"query": q,
	}
	if param != nil {
		qx["variables"] = param
	}
	qJson, err := json.Marshal(qx)
	if err != nil {
		return nil, err
	}
	log.Println("DoRequire body = ", string(qJson))

	body := bytes.NewBuffer([]byte(qJson))
	rsp, err := c.httpClient.Post(c.hosts[0], "application/json", body)

	return rsp, err
}

func (c *Client) DoProxy(req *http.Request) (*http.Response, error) {
	newReq, err := http.NewRequest(req.Method, client.hosts[0], req.Body)
	if err != nil {
		log.Println("NewRequest error:", err)
		return nil, err
	}
	newReq.Header = req.Header
	return c.httpClient.Do(newReq)
}

func DoProxy(req *http.Request) (*http.Response, error) {
	return client.DoProxy(req)
}

func DoRequire(q string, param interface{}) (*http.Response, error) {
	return client.DoRequire(q, param)
}

func (c *Client) Query(q string, param map[string]interface{}) (*http.Response, error) {
	return nil, nil
}

func (c *Client) Mutation(q string, param map[string]interface{}) (*http.Response, error) {

	qx := map[string]interface{}{
		"query":     q,
		"variables": param,
	}
	qJson, err := json.Marshal(qx)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer([]byte(qJson))
	return c.httpClient.Post(c.hosts[0], "application/json", body)
}

func Query(q string, param map[string]interface{}) (*http.Response, error) {
	return client.Query(q, param)
}

func Mutation(q string, param map[string]interface{}) (*http.Response, error) {
	return client.Mutation(q, param)
}
