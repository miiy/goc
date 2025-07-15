package baidu_aip

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/Baidu-AIP/golang-sdk/aip/censor"
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

// TextCensor
// https://cloud.baidu.com/doc/ANTIPORN/s/Rk3h6xb3i
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

// ImgCensor
// https://cloud.baidu.com/doc/ANTIPORN/s/jk42xep4e
func (b *BaiduAip) ImgCensor(file string, options map[string]interface{}) (*CensorResponse, error) {
	fb, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	base64Str := base64.StdEncoding.EncodeToString(fb)
	resStr := b.client.ImgCensor(base64Str, options)
	log.Println("baidu_api.ImgCensor: " + resStr)
	var result CensorResponse
	err = json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ImgCensorUrl
// https://cloud.baidu.com/doc/ANTIPORN/s/jk42xep4e
func (b *BaiduAip) ImgCensorUrl(imgUrl string, options map[string]interface{}) (*CensorResponse, error) {
	resStr := b.client.ImgCensorUrl(imgUrl, options)
	log.Println("baidu_api.ImgCensorUrl: " + resStr)
	var result CensorResponse
	err := json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VoiceCensor
// https://cloud.baidu.com/doc/ANTIPORN/s/hk928u7bz
// rate 音频采样率[16000] ]
// fmt 音频文件的格式，pcm、wav、amr、m4a，推荐pcm格式
func (b *BaiduAip) VoiceCensor(file string, rate int, fmt string, options map[string]interface{}) (*CensorResponse, error) {
	fb, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	base64Str := base64.StdEncoding.EncodeToString(fb)
	resStr := b.client.VoiceCensor(base64Str, rate, fmt, options)
	log.Println("baidu_api.VoiceCensor: " + resStr)
	var result CensorResponse
	err = json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VoiceCensorUrl
// https://cloud.baidu.com/doc/ANTIPORN/s/hk928u7bz
func (b *BaiduAip) VoiceCensorUrl(voiceUrl string, rate int, fmt string, options map[string]interface{}) (*CensorResponse, error) {
	resStr := b.client.VoiceCensorUrl(voiceUrl, rate, fmt, options)
	log.Println("baidu_api.VoiceCensorUrl: " + resStr)
	var result CensorResponse
	err := json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
