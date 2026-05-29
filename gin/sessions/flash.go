package sessions

import (
	"encoding/gob"

	"github.com/gin-gonic/gin"
)

const (
	FlashLevelSuccess = "success"
	FlashLevelError   = "error"
	FlashLevelWarning = "warning"
	FlashLevelInfo    = "info"
)

// flashSessionKey is the session key used to store flash messages.
const flashSessionKey = "_flashes"

// Flash represents a flash message.
type Flash struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type flashSession interface {
	Flashes(vars ...string) []interface{}
}

func init() {
	gob.Register(Flash{})
}

// AddFlash stores a flash message in the session.
func AddFlash(c *gin.Context, flashLevel, message string) error {
	session := Default(c)
	session.AddFlash(Flash{Level: flashLevel, Message: message}, flashSessionKey)
	return session.Save()
}

// Flashes returns and clears flash messages from the session.
func Flashes(c *gin.Context) ([]Flash, error) {
	session := Default(c)
	flashes := flashes(session)
	return flashes, session.Save()
}

func flashes(session flashSession) []Flash {
	values := session.Flashes(flashSessionKey)
	if len(values) == 0 {
		return nil
	}

	flashes := make([]Flash, 0, len(values))
	for _, value := range values {
		flash, ok := value.(Flash)
		if !ok {
			continue
		}
		flashes = append(flashes, flash)
	}
	return flashes
}
