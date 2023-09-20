package tencentcloud

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ims "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ims/v20201229"
)

type IMSClient struct {
	client *ims.Client
}

// ImageModerationRequest
// FileContent 与 FileUrl 二选一
type ImageModerationRequest struct {
	FileContent string
	FileUrl     string
}

func NewIMSClient(credential *common.Credential, region string) (*IMSClient, error) {
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ims.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ims.NewClient(credential, region, cpf)

	return &IMSClient{client: client}, nil
}

// ImageModeration 图片内容安全
// https://cloud.tencent.com/document/api/1125/53273
func (c *IMSClient) ImageModeration(req *ImageModerationRequest) (*ims.ImageModerationResponse, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ims.NewImageModerationRequest()

	if req.FileContent != "" {
		request.FileContent = common.StringPtr(req.FileContent)
	}
	if req.FileUrl != "" {
		request.FileUrl = common.StringPtr(req.FileUrl)
	}

	// 返回的resp是一个ImageModerationResponse的实例，与请求对象对应
	response, err := c.client.ImageModeration(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("An API error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	//// 输出json格式的字符串回包
	//fmt.Printf("%s", response.ToJsonString())
	return response, nil
}
