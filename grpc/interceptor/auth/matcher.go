package auth

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
)

func MatchFullMethods(methods ...string) selector.Matcher {
	if len(methods) == 0 {
		return nil
	}
	methodSet := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		methodSet[method] = struct{}{}
	}
	return selector.MatchFunc(func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		_, ok := methodSet[callMeta.FullMethod()]
		return ok
	})
}
