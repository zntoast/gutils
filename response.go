package gutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, nil, "操作成功", c)
}

func OkWithMessage(msg string, c *gin.Context) {
	Result(SUCCESS, nil, msg, c)
}

func OkWithData(data any, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data any, msg string, c *gin.Context) {
	Result(SUCCESS, data, msg, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, nil, "操作失败", c)
}

func FailWithmessage(msg string, c *gin.Context) {
	Result(ERROR, nil, msg, c)
}

func FailWithDetailed(data any, msg string, c *gin.Context) {
	Result(ERROR, data, msg, c)
}
