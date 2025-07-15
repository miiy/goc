package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/marshal_httpbodyproto.go
type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// A custom marshaler implementation, that doesn't implement the delimited interface
type CustomMarshaler struct {
	runtime.Marshaler
}

func (h *CustomMarshaler) ContentType(v interface{}) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType(v)
}

func (c *CustomMarshaler) Marshal(v interface{}) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetData(), nil
	}
	m, ok := v.(*status.Status)
	if ok {
		return c.Marshaler.Marshal(&Response{
			Code: m.GetCode(),
			Msg:  m.GetMessage(),
			Data: m.GetDetails(),
		})
	}
	return c.Marshaler.Marshal(&Response{
		Code: int32(codes.OK),
		Msg:  "success",
		Data: v,
	})
}
