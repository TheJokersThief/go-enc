package main

import (
	"errors"
)

type ENCNodegroup struct {
	Parent      string                 `json:"parent", yaml:"parent"`
	Classes     map[string]interface{} `json:"classes", yaml:"classes"`
	Nodes       []string               `json:"nodes", yaml:"nodes"`
	Parameters  map[string]interface{} `json:"parameter", yaml:"parameter"`
	Environment string                 `json:"environment", yaml:"environment"`
}

type ENC struct {
	Nodegroups map[string]ENCNodegroup `json:"nodes", yaml:"nodes"`
	ConfigType string
}

func NewENC(config_type string) *ENC {
	return &ENC{
		Nodegroups: map[string]ENCNodegroup{},
		ConfigType: config_type,
	}
}

func (enc *ENC) AddNodegroup(name string, parent string, classes map[string]interface{}, nodes []string, params map[string]interface{}) (*ENCNodegroup, error) {
	nodegroup := ENCNodegroup{
		Parent:     parent,
		Classes:    classes,
		Nodes:      nodes,
		Parameters: params,
	}

	if _, ok := enc.Nodegroups[name]; !ok {
		enc.Nodegroups[name] = nodegroup
	} else {
		return &ENCNodegroup{}, errors.New("Nodegroup already exists")
	}

	return &nodegroup, nil
}

func (enc *ENC) RemoveNodegroup(name string) (*ENCNodegroup, error) {
	if _, ok := enc.Nodegroups[name]; !ok {
		return &ENCNodegroup{}, errors.New("Nodegroup does not exist")
	} else {
		nodegroup := enc.Nodegroups[name]
		delete(enc.Nodegroups, name)
		return &nodegroup, nil
	}
}

func (enc *ENC) GetNodegroup()    {}
func (enc *ENC) AddNode()         {}
func (enc *ENC) RemoveNode()      {}
func (enc *ENC) GetNode()         {}
func (enc *ENC) AddParameter()    {}
func (enc *ENC) SetParameter()    {}
func (enc *ENC) RemoveParameter() {}
func (enc *ENC) AddClass()        {}
func (enc *ENC) RemoveClass()     {}
func (enc *ENC) SetClasses()      {}
func (enc *ENC) SetParent()       {}
