package baidu_translate

import (
	"testing"
)

var testAppId = ""
var testSecKey = ""

func TestBaiduTranslate_Translate(t *testing.T) {
	c := NewBaiduTranslate(testAppId, testSecKey)
	ret, err := c.Translate("苹果", "", EN)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
