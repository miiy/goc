package baidu_aip

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Baidu-AIP/golang-sdk/aip/censor"
	"log"
	"os"
)

// https://cloud.baidu.com/doc/ANTIPORN/s/2kvuvd2pr

type BaiduAip struct {
	client *censor.ContentCensorClient
}

// ConclusionType 1.合规，2.不合规，3.疑似，4.审核失败
type ConclusionType int

const (
	ConclusionTypeOK      ConclusionType = 1
	ConclusionTypeNotOK   ConclusionType = 2
	ConclusionTypePerhaps ConclusionType = 3
	ConclusionTypeFail    ConclusionType = 4
)

type CensorResponse struct {
	LogId          int            `json:"log_id"`
	Conclusion     string         `json:"conclusion"`
	ConclusionType ConclusionType `json:"conclusionType"`
	Data           []interface{}  `json:"data"`
	ErrCode        int            `json:"err_code"`
	ErrMessage     int            `json:"err_message"`
}

func NewClient(apiKey, secretKey string) *BaiduAip {
	client := censor.NewClient(apiKey, secretKey)
	return &BaiduAip{
		client: client,
	}
}

func (b *BaiduAip) TextCensor(text string) (*CensorResponse, error) {
	resStr := b.client.TextCensor(text)
	log.Println("baidu_api.TextCensor: " + resStr)
	var result CensorResponse
	err := json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *BaiduAip) ImgCensor(file string) (*CensorResponse, error) {
	fb, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	base64Str := base64.StdEncoding.EncodeToString(fb)
	resStr := b.client.ImgCensor(base64Str, nil)
	log.Println("baidu_api.ImgCensor: " + resStr)
	var result CensorResponse
	err = json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *BaiduAip) ImgCensorUrl(imgUrl string) (*CensorResponse, error) {
	resStr := b.client.ImgCensorUrl(imgUrl, nil)
	log.Println("baidu_api.ImgCensorUrl: " + resStr)
	var result CensorResponse
	err := json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
