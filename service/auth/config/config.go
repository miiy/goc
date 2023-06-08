package config

import (
	"github.com/miiy/goc/auth/jwt"
)

type Config struct {
	App      App         `yaml:"app"`
	Database Database    `yaml:"database"`
	Server   Server      `yaml:"server"`
	Jwt      jwt.Options `yaml:"jwt"`
}

type App struct {
	Name  string `yaml:"name"`
	Env   string `yaml:"env"`
	Debug bool   `yaml:"debug"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Server struct {
	Http ServerHttp `yaml:"http"`
	Grpc ServerGrpc `yaml:"grpc"`
}

type ServerHttp struct {
	Addr string `yaml:"addr"`
}

type ServerGrpc struct {
	Addr string `yaml:"addr"`
}
