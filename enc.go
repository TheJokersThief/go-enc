package main

import ()

type ENC struct {
	ENCNodes *[]EncNode `json:"nodes"`
}

type ENCNode struct {
	Parent     string             `json:"parent"`
	Classes    *[]interface{}     `json:"classes"`
	Nodes      *[]string          `json:"nodes"`
	Parameters *map[string]string `json:"parameter"`
}
