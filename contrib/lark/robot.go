package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Robot struct {
	hookId string
	client *http.Client
}

type Option func(*Robot)

const hookUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/"

func NewRobot(hookId string) (*Robot, error) {
	return &Robot{
		hookId: hookId,
		client: http.DefaultClient,
	}, nil
}

func WithClient(c *http.Client) Option {
	return func(r *Robot) {
		r.client = c
	}
}

func (r *Robot) SendText(ctx context.Context, msg *TextMessage) error {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, hookUrl+r.hookId, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("send text failed")
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var sResp SendResponse
	if err = json.Unmarshal(respBody, &sResp); err != nil {
		return err
	}
	if sResp.Code != 0 {
		return errors.New(sResp.Msg)
	}
	return nil
}
