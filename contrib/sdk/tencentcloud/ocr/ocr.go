package ocr

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

type Client struct {
	client *ocr.Client
}

func NewClient(credential *common.Credential, region string) (*Client, error) {
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ocr.NewClient(credential, region, cpf)

	return &Client{client: client}, nil
}

func NewGeneralBasicOCRRequest() *ocr.GeneralBasicOCRRequest {
	return ocr.NewGeneralBasicOCRRequest()
}

// GeneralBasicOCR 通用印刷体识别
// https://cloud.tencent.com/document/api/866/33526
func (c *Client) GeneralBasicOCR(request *ocr.GeneralBasicOCRRequest) (*ocr.GeneralBasicOCRResponse, error) {
	// 返回的resp是一个GeneralBasicOCRResponse的实例，与请求对象对应
	response, err := c.client.GeneralBasicOCR(request)
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

func NewGeneralAccurateOCRRequest() *ocr.GeneralAccurateOCRRequest {
	return ocr.NewGeneralAccurateOCRRequest()
}

// GeneralAccurateOCR 通用印刷体识别（高精度版）
// https://cloud.tencent.com/document/api/866/34937
func (c *Client) GeneralAccurateOCR(request *ocr.GeneralAccurateOCRRequest) (*ocr.GeneralAccurateOCRResponse, error) {
	// 返回的resp是一个GeneralAccurateOCRResponse的实例，与请求对象对应
	response, err := c.client.GeneralAccurateOCR(request)
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

func NewIDCardOCRRequest() *ocr.IDCardOCRRequest {
	return ocr.NewIDCardOCRRequest()
}

// IDCardOCR 身份证识别
// https://cloud.tencent.com/document/api/866/33524
func (c *Client) IDCardOCR(request *ocr.IDCardOCRRequest) (*ocr.IDCardOCRResponse, error) {
	// 返回的resp是一个IDCardOCRResponse的实例，与请求对象对应
	response, err := c.client.IDCardOCR(request)
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

func NewBankCardOCRRequest() *ocr.BankCardOCRRequest {
	return ocr.NewBankCardOCRRequest()
}

// BankCardOCR 银行卡识别
// https://cloud.tencent.com/document/api/866/36216
func (c *Client) BankCardOCR(request *ocr.BankCardOCRRequest) (*ocr.BankCardOCRResponse, error) {
	// 返回的resp是一个BankCardOCRResponse的实例，与请求对象对应
	response, err := c.client.BankCardOCR(request)
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

func NewBizLicenseOCRRequest() *ocr.BizLicenseOCRRequest {
	return ocr.NewBizLicenseOCRRequest()
}

// BizLicenseOCR 营业执照识别
// https://cloud.tencent.com/document/api/866/36215
func (c *Client) BizLicenseOCR(request *ocr.BizLicenseOCRRequest) (*ocr.BizLicenseOCRResponse, error) {
	// 返回的resp是一个BizLicenseOCRResponse的实例，与请求对象对应
	response, err := c.client.BizLicenseOCR(request)
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
