package sessions

import (
	"errors"
	"net/http/httptest"
	"testing"

	ginsessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestAddFlashStoresAndSavesFlash(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ginsessions.DefaultKey, session)

	if err := AddFlash(c, FlashLevelSuccess, "saved"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	values := session.flashes[flashSessionKey]
	if len(values) != 1 {
		t.Fatalf("expected 1 flash, got %d", len(values))
	}

	flash, ok := values[0].(Flash)
	if !ok {
		t.Fatalf("expected Flash, got %T", values[0])
	}
	if flash.Level != FlashLevelSuccess || flash.Message != "saved" {
		t.Fatalf("unexpected flash: %#v", flash)
	}
	if !session.saved {
		t.Fatal("expected session to be saved")
	}
}

func TestAddFlashReturnsSaveError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	session.saveErr = errors.New("save failed")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ginsessions.DefaultKey, session)

	if err := AddFlash(c, FlashLevelError, "failed"); err == nil {
		t.Fatal("expected save error")
	}
}

func TestFlashesReturnsAndClearsFlashes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	flash := Flash{Level: FlashLevelInfo, Message: "loaded"}
	session.AddFlash(flash, flashSessionKey)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ginsessions.DefaultKey, session)

	flashes, err := Flashes(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(flashes) != 1 || flashes[0] != flash {
		t.Fatalf("unexpected flashes: %#v", flashes)
	}
	if len(session.flashes[flashSessionKey]) != 0 {
		t.Fatal("expected flashes to be cleared")
	}
	if !session.saved {
		t.Fatal("expected session to be saved")
	}
}

func TestFlashesIgnoresUnexpectedValues(t *testing.T) {
	session := newTestSession()
	flash := Flash{Level: FlashLevelWarning, Message: "kept"}
	session.flashes[flashSessionKey] = []interface{}{"ignored", flash}

	flashes := flashes(session)

	if len(flashes) != 1 || flashes[0] != flash {
		t.Fatalf("unexpected flashes: %#v", flashes)
	}
}
