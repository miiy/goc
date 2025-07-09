package lark

import (
	"context"
	"testing"
)

func TestRobot_SendText(t *testing.T) {
	c, err := NewRobot("")
	if err != nil {
		t.Log(err)
	}
	msg, err := NewTextMsg("test")
	if err != nil {
		t.Log(err)
	}
	if err = c.SendText(context.Background(), msg); err != nil {
		t.Log(err)
	}
}
