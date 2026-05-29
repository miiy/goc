package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type Response struct {
	Code int32 `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type CustomMarshaler struct {
	runtime.Marshaler
}

func (h *CustomMarshaler) ContentType(v any) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType(v)
}

func (c *CustomMarshaler) Marshal(v any) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetData(), nil
	}
	if s, ok := v.(*status.Status); ok {
		return c.Marshaler.Marshal(&Response{
			Code: s.GetCode(),
			Msg:  s.GetMessage(),
			Data: s.GetDetails(),
		})
	}
	return c.Marshaler.Marshal(v)
}
