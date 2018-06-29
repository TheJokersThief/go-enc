package enc

import (
	"errors"
	"reflect"
	"strings"

	"github.com/derekparker/trie"
)

// Nodegroup represents groups of nodes and meta information about them
type Nodegroup struct {
	Parent      string                 `json:"parent" yaml:"parent"`
	Classes     map[string]interface{} `json:"classes" yaml:"classes"`
	Nodes       []string               `json:"nodes" yaml:"nodes"`
	Parameters  map[string]interface{} `json:"parameters" yaml:"parameters"`
	Environment string                 `json:"environment" yaml:"environment"`
}

// ENC represents the entire structure of the External Node Classifier
type ENC struct {
	Nodegroups map[string]Nodegroup
	Nodes      *trie.Trie
	ConfigType string
	FileName   string
}

// NewENC initialises a new ENC
func NewENC(configType string, fileName string) *ENC {
	return &ENC{
		Nodegroups: map[string]Nodegroup{},
		Nodes:      trie.New(),
		ConfigType: configType,
		FileName:   fileName,
	}
}

// AddNodegroup adds a Nodegroup to the ENC
func (enc *ENC) AddNodegroup(name string, parent string, classes map[string]interface{}, nodes []string, params map[string]interface{}) (*Nodegroup, error) {

	nodegroup := Nodegroup{
		Parent:     parent,
		Classes:    classes,
		Parameters: params,
	}

	if len(nodes) == 0 {
		nodegroup.Nodes = []string{}
	}

	if _, ok := enc.Nodegroups[name]; !ok {
		enc.Nodegroups[name] = nodegroup
	} else {
		return &Nodegroup{}, errors.New("Nodegroup already exists")
	}

	if len(nodes) != 0 {
		enc.AddNodes(name, nodes)
	}

	return &nodegroup, nil
}

// RemoveNodegroup removes a nodegroup from the ENC
func (enc *ENC) RemoveNodegroup(name string) (*Nodegroup, error) {
	if val, ok := enc.Nodegroups[name]; ok {
		nodegroup := val
		delete(enc.Nodegroups, name)
		return &nodegroup, nil
	}

	return &Nodegroup{}, errors.New("Nodegroup does not exist")
}

// GetNodegroup retrieves a nodegroup by name
func (enc *ENC) GetNodegroup(nodegroupName string) (*Nodegroup, error) {
	if val, ok := enc.Nodegroups[nodegroupName]; ok {
		return &val, nil
	}

	return &Nodegroup{}, errors.New("Nodegroup does not exist")
}

// AddNode adds a single node to a nodegroup
func (enc *ENC) AddNode(nodegroup string, nodeName string) (*Nodegroup, error) {
	if _, ok := enc.Nodegroups[nodegroup]; !ok {
		return &Nodegroup{}, errors.New("Nodegroup does not exist")
	}
	parentChain := nodeName + "-" + strings.Join(enc.getParentChain(nodegroup), "-")
	if len(parentChain) > len(enc.getLongestChain(nodeName)) {
		enc.Nodes.Add(parentChain, parentChain)
	}

	nodegroupObj, _ := enc.GetNodegroup(nodegroup)
	nodegroupObj.Nodes = append(nodegroupObj.Nodes, nodeName)
	enc.Nodegroups[nodegroup] = *nodegroupObj
	return nodegroupObj, nil
}

// AddNodes adds a slice of nodes to a nodegroup
func (enc *ENC) AddNodes(nodegroup string, nodes []string) (*Nodegroup, error) {
	for _, node := range nodes {
		_, err := enc.AddNode(nodegroup, node)
		if err != nil {
			return &Nodegroup{}, err
		}
	}

	return enc.GetNodegroup(nodegroup)
}

// getParentChain generates a path of parents until it reaches the top
func (enc *ENC) getParentChain(nodegroupName string) []string {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	errCheck(err)

	parents := []string{nodegroupName}
	parent := nodegroup.Parent
	for parent != "" {
		parents = append(parents, parent)
		parentNG, _ := enc.GetNodegroup(parent)
		parent = parentNG.Parent
	}

	return parents
}

// getLongestChain retrieves the current longest parent chain for a node
func (enc *ENC) getLongestChain(node string) string {
	chains := enc.Nodes.PrefixSearch(node)
	longest := ""
	for _, chain := range chains {
		if len(chain) > len(longest) {
			longest = chain
		}
	}

	return longest
}

// RemoveNode removes a single node from a nodegroup
func (enc *ENC) RemoveNode(nodegroup string, nodeName string) (*Nodegroup, error) {
	if _, ok := enc.Nodegroups[nodegroup]; !ok {
		return &Nodegroup{}, errors.New("Nodegroup does not exist")
	}

	chains := enc.Nodes.PrefixSearch(nodeName)
	for _, chain := range chains {

		if strings.HasPrefix(chain, nodeName+"-"+nodegroup) {
			enc.Nodes.Remove(chain)
			break
		}
	}

	nodegroupObj, _ := enc.GetNodegroup(nodegroup)
	nodegroupObj.Nodes = removeByValueSS(nodegroupObj.Nodes, nodeName)
	return nodegroupObj, nil
}

// GetNode retrieves a nodegroup that represents all inherited values for a node
func (enc *ENC) GetNode(nodeName string) (*Nodegroup, error) {
	chain := strings.Replace(enc.getLongestChain(nodeName), nodeName+"-", "", -1)
	path := strings.Split(chain, "-")

	masterNodegroup := &Nodegroup{}

	for _, piece := range reverse(path) {
		pieceNodegroup, _ := enc.GetNodegroup(piece)
		masterNodegroup = enc.mergeNodegroups(masterNodegroup, pieceNodegroup)
	}

	return masterNodegroup, nil
}

// mergeNodegroups merges two nodegroups, preserving values exclusive to ngA and overwriting with
// values from ngB
func (enc *ENC) mergeNodegroups(ngA *Nodegroup, ngB *Nodegroup) *Nodegroup {
	newNG := Nodegroup{
		Parent:      ngB.Parent,
		Nodes:       ngB.Nodes,
		Environment: ngA.Environment,
	}

	if ngB.Environment != "" {
		newNG.Environment = ngB.Environment
	}

	newNG.Classes = enc.mergeNestedMaps(ngA.Classes, ngB.Classes)
	newNG.Parameters = enc.mergeNestedMaps(ngA.Parameters, ngB.Parameters)

	return &newNG
}

// mergeNestedMaps travels nested maps and adds the values from mapB into mapA (overwriting
// what exists and preserving what's unique in both)
func (enc *ENC) mergeNestedMaps(mapA map[string]interface{}, mapB map[string]interface{}) map[string]interface{} {
	for key, val := range mapB {
		if reflect.ValueOf(val).Kind() == reflect.Map {
			if aVal, ok := mapA[key]; ok && reflect.ValueOf(aVal).Kind() == reflect.Map {
				mapA[key] = enc.mergeNestedMaps(aVal.(map[string]interface{}), val.(map[string]interface{}))
			} else {
				if mapA == nil {
					mapA = make(map[string]interface{})
				}
				mapA[key] = val
			}
		} else {
			if mapA == nil {
				mapA = make(map[string]interface{})
			}
			mapA[key] = val
		}
	}

	return mapA
}

// AddParameter add a parameter to a nodegroup
func (enc *ENC) AddParameter(nodegroupName string, key string, val interface{}) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	nodegroup.Parameters[key] = val
	enc.Nodegroups[nodegroupName] = *nodegroup
	return nodegroup, nil
}

// SetParameter is an alias for AddParameter
func (enc *ENC) SetParameter(nodegroupName string, key string, val interface{}) (*Nodegroup, error) {
	return enc.AddParameter(nodegroupName, key, val)
}

// RemoveParameter removes a parameter from a nodegroup
func (enc *ENC) RemoveParameter(nodegroupName string, key string) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	if _, ok := nodegroup.Parameters[key]; !ok {
		return &Nodegroup{}, errors.New("That parameter does not exist for this nodegroup")
	}

	delete(nodegroup.Parameters, key)
	return nodegroup, nil
}

// AddClass initialises a new class for a nodegroup
func (enc *ENC) AddClass(nodegroupName string, key string) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	nodegroup.Classes[key] = make(map[string]interface{})
	enc.Nodegroups[nodegroupName] = *nodegroup
	return nodegroup, nil
}

// RemoveClass removes a class from a nodegroup
func (enc *ENC) RemoveClass(nodegroupName string, key string) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	if _, ok := nodegroup.Classes[key]; !ok {
		return &Nodegroup{}, errors.New("That parameter does not exist for this nodegroup")
	}

	delete(nodegroup.Classes, key)
	return nodegroup, nil
}

// AddClassParameter adds a parameter to a given class on a nodegroup
func (enc *ENC) AddClassParameter(nodegroupName string, class string, key string, val interface{}) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	if _, ok := nodegroup.Classes[class]; !ok {
		return &Nodegroup{}, errors.New("Class does not exist on that nodegroup")
	}

	ngClass := nodegroup.Classes[class].(map[string]interface{})
	ngClass[key] = val
	nodegroup.Classes[class] = ngClass

	return nodegroup, nil
}

// RemoveClassParameter removes a parameter from a class
func (enc *ENC) RemoveClassParameter(nodegroupName string, class string, key string) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	if _, ok := nodegroup.Classes[class]; !ok {
		return &Nodegroup{}, errors.New("Class does not exist on that nodegroup")
	}

	ngClass := nodegroup.Classes[class].(map[string]interface{})
	if _, ok := ngClass[key]; !ok {
		return &Nodegroup{}, errors.New("Parameter for that class does not exist on that nodegroup")
	}

	delete(ngClass, key)
	nodegroup.Classes[class] = ngClass

	return nodegroup, nil
}

// SetClassParameter is an alias for AddClassParameter
func (enc *ENC) SetClassParameter(nodegroupName string, class string, key string, val interface{}) (*Nodegroup, error) {
	return enc.AddClassParameter(nodegroupName, class, key, val)
}

// SetParent sets the value of the parent of a nodegroup
func (enc *ENC) SetParent(nodegroupName string, parent string) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	_, parentErr := enc.GetNodegroup(parent)
	if parentErr != nil {
		return &Nodegroup{}, errors.New("The parent nodegroup does not exist")
	}

	nodegroup.Parent = parent
	enc.Nodegroups[nodegroupName] = *nodegroup
	return nodegroup, nil
}

// SetEnvironment sets the value of the environment of a nodegroup
func (enc *ENC) SetEnvironment(nodegroupName string, env string) (*Nodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	if err != nil {
		return &Nodegroup{}, err
	}

	nodegroup.Environment = env
	enc.Nodegroups[nodegroupName] = *nodegroup
	return nodegroup, nil
}
