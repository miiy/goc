package ims

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ims "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ims/v20201229"
)

type Client struct {
	client *ims.Client
}

func NewClient(credential *common.Credential, region string) (*Client, error) {
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ims.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ims.NewClient(credential, region, cpf)

	return &Client{client: client}, nil
}

func NewImageModerationRequest() *ims.ImageModerationRequest {
	return ims.NewImageModerationRequest()
}

// ImageModeration 图片内容安全
// https://cloud.tencent.com/document/api/1125/53273
func (c *Client) ImageModeration(request *ims.ImageModerationRequest) (*ims.ImageModerationResponse, error) {
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
