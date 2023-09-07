package strings

import "testing"

func TestMaskIdCard(t *testing.T) {
	s := MaskIdCard("100000199901021234")
	t.Log(s)
}

func TestMaskPhone(t *testing.T) {
	s := MaskPhone("13010001234")
	t.Log(s)
}
