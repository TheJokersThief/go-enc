package main

import ()

type Config struct {
	ENC *ENC `json:"enc"`
}

func NewConfig() *Config {
	return &Config{}
}
