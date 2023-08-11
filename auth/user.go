package auth

import (
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/db/gorm"
)

type User struct {
	gorm.Model
	Username          string
	Password          string
	Email             string
	EmailVerifiedTime *db.JSONTime
	Phone             string
	Unionid           string
	MpOpenid          string
	MpSessionKey      string
	Status            int64
}
