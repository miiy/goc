package server

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestWithAuthRequiresApplicationRPCsByDefault(t *testing.T) {
	config := WithAuth(nil, RequireAuthExcept())

	if !matchesRequiredAuth(config, "/nova.test.v1.Service/Protected") {
		t.Fatal("expected application method to require auth")
	}
	if matchesRequiredAuth(config, healthpb.Health_Check_FullMethodName) {
		t.Fatal("expected health check to stay anonymous")
	}
}

func TestWithAuthAllowsAnonymousMethods(t *testing.T) {
	config := WithAuth(nil, RequireAuthExcept("/nova.test.v1.Service/Public"))

	if matchesRequiredAuth(config, "/nova.test.v1.Service/Public") {
		t.Fatal("expected anonymous method to skip required auth")
	}
	if !matchesRequiredAuth(config, "/nova.test.v1.Service/Protected") {
		t.Fatal("expected unlisted application method to require auth")
	}
	if matchesRequiredAuth(config, healthpb.Health_Check_FullMethodName) {
		t.Fatal("expected health check to stay anonymous")
	}
}

func matchesRequiredAuth(config *AuthConfig, method string) bool {
	if config == nil || config.requireAuth == nil {
		return false
	}
	return config.requireAuth.Match(context.Background(), interceptors.NewServerCallMeta(method, nil, nil))
}
