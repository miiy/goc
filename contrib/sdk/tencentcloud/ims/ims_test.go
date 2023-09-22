package ims

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

func TestImageModeration(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
	}
	req := NewImageModerationRequest()
	url := ""
	req.FileUrl = &url
	resp, err := client.ImageModeration(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}
