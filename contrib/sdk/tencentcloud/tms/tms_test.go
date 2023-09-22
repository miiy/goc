package tms

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

func TestTextModeration(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
	}
	resp, err := client.TextModeration("你好")
	if err != nil {
		t.Error(err)
	}
	if *resp.Response.Suggestion == "Block" || *resp.Response.Suggestion == "Review" {
		t.Error("Block or Review")
	}
	t.Log(resp.ToJsonString())
}
