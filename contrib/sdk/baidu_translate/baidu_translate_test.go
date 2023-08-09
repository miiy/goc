package baidu_translate

import (
	"testing"
)

var appId = ""
var secKey =""

func TestBaiduTranslate_Translate(t *testing.T) {
	c := NewBaiduTranslate(appId, secKey)
	ret, err := c.Translate("苹果", "", EN)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
