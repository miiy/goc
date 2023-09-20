package tencentcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	"testing"
)

func TestTextModeration(t *testing.T) {
	client, err := NewTmsClient(testCredential, regions.Beijing)
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
