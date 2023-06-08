package dingtalk_robot

type TextMsg struct {
	At      TextMsgAt   `json:"at"`
	Text    TextMsgText `json:"text"`
	MsgType string      `json:"msgtype"`
}
type TextMsgAt struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type TextMsgText struct {
	Content string `json:"content"`
}

type SendResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewTextMsg(content string) *TextMsg {
	return &TextMsg{
		At: TextMsgAt{
			AtMobiles: []string{},
			AtUserIds: []string{},
		},
		Text: TextMsgText{
			Content: content,
		},
		MsgType: "text",
	}
}
