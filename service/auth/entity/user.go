package entity

import (
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/db/gorm"
)

const (
	UserColumnUsername = "username"
	UserColumnEmail    = "email"
	UserColumnPhone    = "phone"
)

type User struct {
	gorm.Model
	Username          string
	Password          string
	Email             string
	EmailVerifiedTime *db.JSONTime
	Phone             string
	Status            int64
}
