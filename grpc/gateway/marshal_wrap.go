package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/genproto/googleapis/rpc/status"
)

// WrapResponse wraps all responses (success and error) in a unified format.
//
// Error response:
//
//	{"code": 6, "message": "...", "data": null}
//
// Success response:
//
//	{"code": 0, "message": "ok", "data": {...}}
type WrapResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// WrapMarshaler wraps all responses in a unified {code, message, data} format.
type WrapMarshaler struct {
	runtime.Marshaler
}

func (h *WrapMarshaler) ContentType(v any) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType(v)
}

func (c *WrapMarshaler) Marshal(v any) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetData(), nil
	}
	if s, ok := v.(*status.Status); ok {
		return c.Marshaler.Marshal(&WrapResponse{
			Code:    s.GetCode(),
			Message: s.GetMessage(),
			Data:    nil,
		})
	}
	return c.Marshaler.Marshal(&WrapResponse{
		Code:    0,
		Message: "ok",
		Data:    v,
	})
}
