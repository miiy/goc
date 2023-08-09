package baidu_aip

import (
	"encoding/json"
	"github.com/Baidu-AIP/golang-sdk/aip/censor"
	"log"
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

type TextCensorResponse struct {
	LogId          int                      `json:"log_id"`
	Conclusion     string                   `json:"conclusion"`
	ConclusionType ConclusionType           `json:"conclusionType"`
	Data           []TextCensorResponseData `json:"data"`
	ErrCode        int                      `json:"err_code"`
	ErrMessage     int                      `json:"err_message"`
}

type TextCensorResponseData struct {
	Type           int                          `json:"type"`
	SubType        int                          `json:"subType"`
	Conclusion     string                       `json:"conclusion"`
	ConclusionType ConclusionType               `json:"conclusionType"`
	Msg            string                       `json:"msg"`
	Hits           []TextCensorResponseDataHits `json:"hits"`
}

type TextCensorResponseDataHits struct {
	DataSetName string
	Probability float64
	Words       []string
}

func NewClient(apiKey, secretKey string) *BaiduAip {
	client := censor.NewClient(apiKey, secretKey)
	return &BaiduAip{
		client: client,
	}
}

func (b *BaiduAip) TextCensor(text string) (*TextCensorResponse, error) {
	resultStr := b.client.TextCensor(text)
	log.Println("baidu_api.TextCensor: " + resultStr)
	var result TextCensorResponse
	err := json.Unmarshal([]byte(resultStr), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
