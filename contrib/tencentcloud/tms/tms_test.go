package tms

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
