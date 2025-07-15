package baidu_aip

import "testing"

const (
	testApiKey    = ""
	testSecretKey = ""
)

var tc *BaiduAip

func TestNewClient(t *testing.T) {
	tc = NewClient(testApiKey, testSecretKey)
	t.Logf("%+v", tc)
}

func TestTextCensor(t *testing.T) {
	TestNewClient(t)
	ret, err := tc.TextCensor("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", ret)
	if ret.ConclusionType == ConclusionTypeOK {
		t.Log("ok")
	}
}

func TestImgCensor(t *testing.T) {
	TestNewClient(t)
	ret, err := tc.ImgCensor("./test.png", nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", ret)
	if ret.ConclusionType == ConclusionTypeOK {
		t.Log("ok")
	}
}

func TestImgCensorUrl(t *testing.T) {
	TestNewClient(t)
	ret, err := tc.ImgCensorUrl("", nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", ret)
	if ret.ConclusionType == ConclusionTypeOK {
		t.Log("ok")
	}
}
