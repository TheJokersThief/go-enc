package main

import ()

type ENC struct {
	ENCNodes *[]EncNode `json:"nodes"`
}

type ENCNode struct {
	Parent     string                  `json:"parent"`
	Classes    *map[string]interface{} `json:"classes"`
	Nodes      *[]string               `json:"nodes"`
	Parameters *map[string]interface{} `json:"parameter"`
}
