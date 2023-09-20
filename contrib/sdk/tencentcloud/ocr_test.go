package tencentcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	"testing"
)

func TestIDCardOCR(t *testing.T) {
	client, err := NewOCRClient(testCredential, regions.Beijing)
	if err != nil {
		t.Error(err)
	}
	req := IDCardOCRRequest{ImageUrl: ""}
	resp, err := client.IDCardOCR(&req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}

func TestBankCardOCR(t *testing.T) {
	client, err := NewOCRClient(testCredential, regions.Beijing)
	if err != nil {
		t.Error(err)
	}
	req := BankCardOCRRequest{ImageUrl: ""}
	resp, err := client.BankCardOCR(&req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}

func TestBizLicenseOCR(t *testing.T) {
	client, err := NewOCRClient(testCredential, regions.Beijing)
	if err != nil {
		t.Error(err)
	}
	req := BizLicenseOCRRequest{ImageUrl: ""}
	resp, err := client.BizLicenseOCR(&req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}
