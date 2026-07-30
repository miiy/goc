package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
)

const MaxUserAgentLength = 512

// UserAgent returns the first User-Agent header value, truncated to the
// supported storage length.
func UserAgent(c *gin.Context) string {
	return truncateRunes(c.GetHeader("User-Agent"), MaxUserAgentLength)
}

func truncateRunes(value string, max int) string {
	if max < 1 {
		return ""
	}
	count := 0
	for index := range value {
		if count == max {
			return value[:index]
		}
		count++
	}
	return value
}

// ExtractToken extracts a Bearer token from the Authorization header.
func ExtractToken(c *gin.Context) (string, error) {
	return request.BearerExtractor{}.ExtractToken(c.Request)
}
