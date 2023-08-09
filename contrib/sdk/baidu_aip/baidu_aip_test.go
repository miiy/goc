package baidu_aip

import "testing"

const (
	apiKey    = ""
	secretKey = ""
)

var tc *BaiduAip

func TestNewClient(t *testing.T) {
	tc = NewClient(apiKey, secretKey)
	t.Logf("%+v", tc)
}

func TestTextCensor(t *testing.T) {
	TestNewClient(t)
	ret, err := tc.TextCensor("测试")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", ret)
	if ret.ConclusionType == ConclusionTypeOK {
		t.Log("ok")
	}
}
