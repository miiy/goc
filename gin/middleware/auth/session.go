package auth

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin/sessions"
)

const SessionKeyAuthUser = "goc.auth"

func SessionUser(value any) (*gocauth.AuthenticatedUser, bool) {
	values, ok := value.(map[string]any)
	if !ok {
		return nil, false
	}
	return sessionUserFromMap(values)
}

func SessionAuthenticationMiddleware(redirectPath string) gin.HandlerFunc {
	if redirectPath == "" {
		redirectPath = "/register"
	}

	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user, ok := SessionUser(session.Get(SessionKeyAuthUser))
		if !ok {
			ctx.Redirect(http.StatusFound, redirectPath)
			ctx.Abort()
			return
		}

		setAuthUser(ctx, user)
		ctx.Next()
	}
}

func sessionUserFromMap(values map[string]any) (*gocauth.AuthenticatedUser, bool) {
	username, _ := values["username"].(string)
	if username == "" {
		return nil, false
	}

	id, ok := sessionUserID(values["id"])
	if !ok {
		return nil, false
	}
	return &gocauth.AuthenticatedUser{ID: id, Username: username}, true
}

func sessionUserID(value any) (int64, bool) {
	switch id := value.(type) {
	case int:
		return int64(id), id > 0
	case int8:
		return int64(id), id > 0
	case int16:
		return int64(id), id > 0
	case int32:
		return int64(id), id > 0
	case int64:
		return id, id > 0
	case uint:
		if uint64(id) > uint64(^uint(0)>>1) {
			return 0, false
		}
		return int64(id), id > 0
	case uint8:
		return int64(id), id > 0
	case uint16:
		return int64(id), id > 0
	case uint32:
		return int64(id), id > 0
	case uint64:
		if id > uint64(^uint(0)>>1) {
			return 0, false
		}
		return int64(id), id > 0
	case float32:
		return sessionUserFloatID(float64(id))
	case float64:
		return sessionUserFloatID(id)
	case json.Number:
		parsed, err := id.Int64()
		return parsed, err == nil && parsed > 0
	case string:
		parsed, err := strconv.ParseInt(id, 10, 64)
		return parsed, err == nil && parsed > 0
	default:
		return 0, false
	}
}

func sessionUserFloatID(id float64) (int64, bool) {
	const maxInt64Exclusive = float64(1 << 63)

	if math.IsNaN(id) || math.IsInf(id, 0) || id <= 0 || math.Trunc(id) != id || id >= maxInt64Exclusive {
		return 0, false
	}
	return int64(id), true
}
