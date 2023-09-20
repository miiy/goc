package tencentcloud

import (
	"encoding/base64"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20201229"
)

type TmsClient struct {
	client *tms.Client
}

func NewTmsClient(c *common.Credential, region string) (*TmsClient, error) {
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tms.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := tms.NewClient(c, region, cpf)
	return &TmsClient{client: client}, nil
}

// TextModeration 文本内容安全
// https://cloud.tencent.com/document/api/1124/51860
func (c *TmsClient) TextModeration(str string) (*tms.TextModerationResponse, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tms.NewTextModerationRequest()

	encodeContent := base64.StdEncoding.EncodeToString([]byte(str))
	request.Content = common.StringPtr(encodeContent)

	// 返回的resp是一个TextModerationResponse的实例，与请求对象对应
	response, err := c.client.TextModeration(request)
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
