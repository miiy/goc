package dingtalk_robot

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
	AccessToken string
	Client      *http.Client
}

type Option func(*Robot)

const (
	RobotHost = "https://oapi.dingtalk.com/robot/send"
)

func NewRobot(accessToken string, opts ...Option) (*Robot, error) {
	c := &Robot{
		AccessToken: accessToken,
		Client:      http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

func WithClient(c *http.Client) Option {
	return func(r *Robot) {
		r.Client = c
	}
}

func (r *Robot) SendText(ctx context.Context, msg *TextMsg) error {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?access_token=%s", RobotHost, r.AccessToken), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("send text failed")
	}
	defer resp.Body.Close()
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
