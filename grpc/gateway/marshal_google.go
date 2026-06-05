package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

// GoogleErrorResponse follows Google HTTP JSON API error format.
// See: https://cloud.google.com/apis/design/errors
//
// Error response:
//
//	{"error": {"code": 6, "message": "...", "details": []}}
//
// Success response: proto JSON as-is.
type GoogleErrorResponse struct {
	Error GoogleErrorStatus `json:"error"`
}

type GoogleErrorStatus struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Status  string        `json:"status,omitempty"`
	Details []*anypb.Any `json:"details,omitempty"`
}

// GoogleMarshaler wraps errors in Google HTTP JSON API format.
type GoogleMarshaler struct {
	runtime.Marshaler
}

func (h *GoogleMarshaler) ContentType(v any) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType(v)
}

func (c *GoogleMarshaler) Marshal(v any) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetData(), nil
	}
	if s, ok := v.(*status.Status); ok {
		return c.Marshaler.Marshal(&GoogleErrorResponse{
			Error: GoogleErrorStatus{
				Code:    s.GetCode(),
				Message: s.GetMessage(),
				Details: s.GetDetails(),
			},
		})
	}
	return c.Marshaler.Marshal(v)
}
