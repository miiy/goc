package ocr

import (
	"github.com/miiy/goc/contrib/sdk/tencentcloud"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	"testing"
)

var (
	testSecretId  = ""
	testSecretKey = ""
)

func newTestClient() (*Client, error) {
	return NewClient(tencentcloud.NewCredential(testSecretId, testSecretKey), regions.Beijing)
}

func TestGeneralBasicOCR(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
	}

	req := NewGeneralBasicOCRRequest()
	url := ""
	req.ImageUrl = &url
	resp, err := client.GeneralBasicOCR(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}

func TestIDCardOCR(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
	}

	req := NewIDCardOCRRequest()
	url := ""
	req.ImageUrl = &url
	resp, err := client.IDCardOCR(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}

func TestBankCardOCR(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
	}

	req := NewBankCardOCRRequest()
	url := ""
	req.ImageUrl = &url
	resp, err := client.BankCardOCR(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}

func TestBizLicenseOCR(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
	}

	req := NewBizLicenseOCRRequest()
	url := ""
	req.ImageUrl = &url
	resp, err := client.BizLicenseOCR(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}
