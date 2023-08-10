package config

import (
	conf "github.com/miiy/goc/config"
)

type Config struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
	Server   Server   `yaml:"server"`
	Jwt      Jwt      `yaml:"jwt"`
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

type Redis struct {
	Addrs    []string `yaml:"addrs"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	DB       int      `yaml:"db"`
}

type Jwt struct {
	Secret    string `yaml:"secret"`
	Issuer    string `yaml:"issuer"`
	ExpiresIn int64  `yaml:"expiresIn"`
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

var config *Config

func NewConfig(fileName string) (*Config, error) {
	if err := conf.Load(fileName, &config); err != nil {
		return nil, err
	}
	return config, nil
}
