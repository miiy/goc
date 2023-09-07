package rand

import "testing"

func TestRandInt(t *testing.T) {
	n := RandInt(100000, 999999)
	t.Log(n)
}
