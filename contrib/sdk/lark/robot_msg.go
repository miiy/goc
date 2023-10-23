package lark

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

type TextMessage struct {
	MsgType   string             `json:"msg_type"`
	Content   TextMessageContent `json:"content"`
	Timestamp string             `json:"timestamp"`
	Sign      string             `json:"sign"`
}

type TextMessageContent struct {
	Text string `json:"text"`
}

type SendResponse struct {
	StatusCode    int         `json:"StatusCode"`
	StatusMessage string      `json:"StatusMessage"`
	Code          int         `json:"code"`
	Data          interface{} `json:"data"`
	Msg           string      `json:"msg"`
}

func NewTextMsg(content string) (*TextMessage, error) {
	// sign
	timestamp := time.Now().Unix()
	sign, err := genSign(content, timestamp)
	if err != nil {
		return nil, err
	}
	msg := TextMessage{
		MsgType: "text",
		Content: TextMessageContent{
			Text: content,
		},
		Timestamp: strconv.Itoa(int(timestamp)),
		Sign:      sign,
	}
	return &msg, nil
}

func genSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
