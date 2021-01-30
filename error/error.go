package error

import (
	"encoding/json"
	"runtime"
)

const (
	SUCCESS        = 0
	ERROR          = 500
	INVALID_PARAMS = 400

	// ERROR_AUTH
	ERROR_NOT_AUTH                 = 20000
	ERROR_AUTH_FAILED              = 20010
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20020
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20030
	ERROR_AUTH_TOKEN               = 20040
	ERROR_AUTH                     = 20050

	// ERROR_DGRAPH
	ERROR_DGRAPH_FAILED = 30000

	// JSON_Mashal_
	ERROR_JSON_MARSHAL_FAILED   = 40000
	ERROR_JSON_UNMARSHAL_FAILED = 40001
	ERROR_USERDATA_ERROR        = 41000
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",
	//ERROR_EXIST_TAG:                 "已存在该标签名称",
	//ERROR_EXIST_TAG_FAIL:            "获取已存在标签失败",
	//ERROR_NOT_EXIST_TAG:             "该标签不存在",
	//ERROR_GET_TAGS_FAIL:             "获取所有标签失败",
	//ERROR_COUNT_TAG_FAIL:            "统计标签失败",
	//ERROR_ADD_TAG_FAIL:              "新增标签失败",
	//ERROR_EDIT_TAG_FAIL:             "修改标签失败",
	//ERROR_DELETE_TAG_FAIL:           "删除标签失败",
	//ERROR_EXPORT_TAG_FAIL:           "导出标签失败",
	//ERROR_IMPORT_TAG_FAIL:           "导入标签失败",
	//ERROR_NOT_EXIST_ARTICLE:         "该文章不存在",
	//ERROR_ADD_ARTICLE_FAIL:          "新增文章失败",
	//ERROR_DELETE_ARTICLE_FAIL:       "删除文章失败",
	//ERROR_CHECK_EXIST_ARTICLE_FAIL:  "检查文章是否存在失败",
	//ERROR_EDIT_ARTICLE_FAIL:         "修改文章失败",
	//ERROR_COUNT_ARTICLE_FAIL:        "统计文章失败",
	//ERROR_GET_ARTICLES_FAIL:         "获取多个文章失败",
	//ERROR_GET_ARTICLE_FAIL:          "获取单个文章失败",
	//ERROR_GEN_ARTICLE_POSTER_FAIL:   "生成文章海报失败",
	ERROR_AUTH_FAILED:              "错误的用户名或密码",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",

	ERROR_DGRAPH_FAILED: "数据库请求失败",
	//ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	//ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	//ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	ERROR_JSON_MARSHAL_FAILED:   "JSON编码失败",
	ERROR_JSON_UNMARSHAL_FAILED: "JSON解码失败",

	ERROR_USERDATA_ERROR: "用户数据出错",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

// it from dgraph
type Location struct {
	Line   int    `json:"line"`
	Column int    `json:"column"`
	File   string `json:"file"`
}

type Error struct {
	Message string      `json:"message"`
	File    string      `json:"file"`
	Line    int         `json:"line"`
	Ctx     interface{} `json:"context"'`
}

func (e *Error) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "format error failed"
	}
	return string(b)
}

func NewError(msg string, ctx interface{}) Error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return Error{
			Message: msg,
		}
	}
	return Error{
		Message: msg,
		File:    file,
		Line:    line,
		Ctx:     ctx,
	}
}
