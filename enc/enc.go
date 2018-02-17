package main

import ()

type ENCNodegroup struct {
	Parent     string                  `json:"parent"`
	Classes    *map[string]interface{} `json:"classes"`
	Nodes      *[]string               `json:"nodes"`
	Parameters *map[string]interface{} `json:"parameter"`
}

type ENC struct {
	ENCNodes *map[string]ENCNodegroup `json:"nodes"`
}

func NewENC() *ENC {
	return &ENC{
		ENCNodes: &map[string]ENCNodegroup{},
	}
}
