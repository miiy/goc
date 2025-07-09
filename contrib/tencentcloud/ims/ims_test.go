package ims

import (
	"testing"

	"github.com/miiy/goc/contrib/tencentcloud"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
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
