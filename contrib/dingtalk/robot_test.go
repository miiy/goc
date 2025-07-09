package dingtalk

import (
	"context"
	"testing"
)

func TestRobot_SendText(t *testing.T) {
	c, err := NewRobot("")
	if err != nil {
		t.Log(err)
	}
	if err = c.SendText(context.Background(), NewTextMsg("test")); err != nil {
		t.Log(err)
	}
}
