package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Robot struct {
	accessToken string
	client      *http.Client
}

type Option func(*Robot)

const (
	HookUrl = "https://oapi.dingtalk.com/robot/send"
)

func NewRobot(accessToken string, opts ...Option) (*Robot, error) {
	c := &Robot{
		accessToken: accessToken,
		client:      http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

func WithClient(c *http.Client) Option {
	return func(r *Robot) {
		r.client = c
	}
}

func (r *Robot) SendText(ctx context.Context, msg *TextMsg) error {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?access_token=%s", HookUrl, r.accessToken), bytes.NewReader(reqBody))
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
	if sResp.ErrCode != 0 {
		return errors.New(sResp.ErrMsg)
	}
	return nil
}
