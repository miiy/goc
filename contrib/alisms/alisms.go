package alisms

import (
	"encoding/json"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

const endpoint = "dysmsapi.aliyuncs.com"

type Config struct {
	AccessKey       string
	AccessKeySecret string
	SignName        string
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{config: config}
}

func (c *Client) CreateClient() (*dysmsapi20170525.Client, error) {
	cred, err := credential.NewCredential(&credential.Config{
		Type:            tea.String("access_key"),
		AccessKeyId:     tea.String(c.config.AccessKey),
		AccessKeySecret: tea.String(c.config.AccessKeySecret),
	})
	if err != nil {
		return nil, err
	}

	config := &openapi.Config{
		Credential: cred,
	}
	config.Endpoint = tea.String(endpoint)
	return dysmsapi20170525.NewClient(config)
}

// SendSms 使用指定模板发送短信，templateParam 为模板变量。
func (c *Client) SendSms(tel, templateCode string, templateParam map[string]string) (string, error) {
	client, err := c.CreateClient()
	if err != nil {
		return "", err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(tel),
		SignName:     tea.String(c.config.SignName),
		TemplateCode: tea.String(templateCode),
	}
	if len(templateParam) > 0 {
		templateParamJSON, err := json.Marshal(templateParam)
		if err != nil {
			return "", err
		}
		sendSmsRequest.TemplateParam = tea.String(string(templateParamJSON))
	}

	resp, err := client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
	if err != nil {
		return "", err
	}
	if resp == nil || resp.Body == nil {
		return "{}", nil
	}
	jsonResult, err := json.Marshal(resp.Body)
	if err != nil {
		return "", err
	}
	return string(jsonResult), nil
}

// SendSmsWithTemplate 使用指定模板 code 发送短信。
func (c *Client) SendSmsWithTemplate(tel string, templateCode string) (string, error) {
	return c.SendSms(tel, templateCode, nil)
}
