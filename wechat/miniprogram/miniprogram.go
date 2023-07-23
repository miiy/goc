package miniprogram

import (
	"context"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	mpConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

type MiniProgram struct {
	MiniProgram *miniprogram.MiniProgram
}

var mp *MiniProgram

func NewMiniProgram(appId, appSecret string) (*MiniProgram, error) {
	memory := cache.NewMemory()
	cfg := &mpConfig.Config{
		AppID:     appId,
		AppSecret: appSecret,
		// EncodingAESKey: "xxxx",
		Cache: memory,
	}
	mimiProgram := miniprogram.NewMiniProgram(cfg)
	mp = &MiniProgram{
		MiniProgram: mimiProgram,
	}
	return mp, nil
}

func (m *MiniProgram) Code2Session(ctx context.Context, code string) (*auth.ResCode2Session, error) {
	ret, err := m.MiniProgram.GetAuth().Code2SessionContext(ctx, code)
	return &ret, err
}
