package tencentcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	"testing"
)

func TestImageModeration(t *testing.T) {
	client, err := NewIMSClient(testCredential, regions.Beijing)
	if err != nil {
		t.Error(err)
	}
	req := ImageModerationRequest{FileUrl: ""}
	resp, err := client.ImageModeration(&req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.ToJsonString())
}
