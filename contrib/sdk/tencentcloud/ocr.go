package tencentcloud

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

type OCRClient struct {
	client *ocr.Client
}

type OCRImageRequest struct {
	ImageBase64 string
	ImageUrl    string
}

type IDCardOCRRequest = OCRImageRequest

type BankCardOCRRequest = OCRImageRequest

type BizLicenseOCRRequest = OCRImageRequest

func NewOCRClient(credential *common.Credential, region string) (*OCRClient, error) {
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ocr.NewClient(credential, region, cpf)

	return &OCRClient{client: client}, nil
}

// IDCardOCR 身份证识别
// https://cloud.tencent.com/document/api/866/33524
func (c *OCRClient) IDCardOCR(req *IDCardOCRRequest) (*ocr.IDCardOCRResponse, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewIDCardOCRRequest()

	if req.ImageUrl != "" {
		request.ImageUrl = common.StringPtr(req.ImageUrl)
	}
	if req.ImageBase64 != "" {
		request.ImageBase64 = common.StringPtr(req.ImageBase64)
	}

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

// BankCardOCR 银行卡识别
// https://cloud.tencent.com/document/api/866/36216
func (c *OCRClient) BankCardOCR(req *BankCardOCRRequest) (*ocr.BankCardOCRResponse, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewBankCardOCRRequest()

	if req.ImageUrl != "" {
		request.ImageUrl = common.StringPtr(req.ImageUrl)
	}
	if req.ImageBase64 != "" {
		request.ImageBase64 = common.StringPtr(req.ImageBase64)
	}

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

// BizLicenseOCR 营业执照识别
// https://cloud.tencent.com/document/api/866/36215
func (c *OCRClient) BizLicenseOCR(req *BizLicenseOCRRequest) (*ocr.BizLicenseOCRResponse, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewBizLicenseOCRRequest()

	if req.ImageUrl != "" {
		request.ImageUrl = common.StringPtr(req.ImageUrl)
	}
	if req.ImageBase64 != "" {
		request.ImageBase64 = common.StringPtr(req.ImageBase64)
	}

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
