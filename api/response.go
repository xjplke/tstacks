package api

import (
	"strings"
	e "techstacks.cn/techstacks/error"
)

type Response struct {
	Code       int                    `json:"code"`
	Error      e.Error                `json:"errors,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func ToString(errs []e.Error) string {
	b := new(strings.Builder)

	for _, e := range errs {
		b.WriteString(e.Error())
	}

	return b.String()
}
