package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CodeSuccess         = 0
	CodeError           = 1
	CodeBadRequest      = 2
	CodeUnauthenticated = 3
)

func CodeText(code int) string {
	switch code {
	case CodeSuccess:
		return "Success"
	case CodeError:
		return "Error"
	case CodeBadRequest:
		return "Bad request"
	case CodeUnauthenticated:
		return "Unauthorized"
	default:
		return ""
	}
}

type Response struct {
	Status int
	Data   ResponseData
}

type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Option func(*Response)

func WithStatus(status int) Option {
	return func(c *Response) {
		c.Status = status
	}
}

func WithCode(code int) Option {
	return func(c *Response) {
		c.Data.Code = code
		c.Data.Msg = CodeText(code)
	}
}

func WithMsg(msg string) Option {
	return func(c *Response) {
		c.Data.Msg = msg
	}
}

func WithData(data any) Option {
	return func(c *Response) {
		c.Data.Data = data
	}
}

func Success(c *gin.Context, v interface{}, opts ...Option) {
	resp := Response{
		Status: http.StatusOK,
		Data: ResponseData{
			Code: CodeSuccess,
			Msg:  CodeText(CodeSuccess),
			Data: v,
		},
	}
	for _, o := range opts {
		o(&resp)
	}

	c.JSON(resp.Status, &resp.Data)
}

func Error(c *gin.Context, opts ...Option) {
	resp := Response{
		Status: http.StatusBadRequest,
		Data: ResponseData{
			Code: CodeError,
			Msg:  CodeText(CodeError),
			Data: nil,
		},
	}
	for _, o := range opts {
		o(&resp)
	}
	c.JSON(resp.Status, &resp.Data)
}
