package enc

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/derekparker/trie"
	// "github.com/jinzhu/copier"
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
	Nodes      *trie.Trie
	ConfigType string
}

func NewENC(config_type string) *ENC {
	return &ENC{
		Nodegroups: map[string]ENCNodegroup{},
		Nodes:      trie.New(),
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
	if val, ok := enc.Nodegroups[name]; ok {
		nodegroup := val
		delete(enc.Nodegroups, name)
		return &nodegroup, nil
	} else {
		return &ENCNodegroup{}, errors.New("Nodegroup does not exist")
	}
}

func (enc *ENC) GetNodegroup(nodegroup_name string) (*ENCNodegroup, error) {
	if val, ok := enc.Nodegroups[nodegroup_name]; ok {
		return &val, nil
	} else {
		return &ENCNodegroup{}, errors.New("Nodegroup does not exist")
	}
}

func (enc *ENC) AddNode(nodegroup string, node_name string) (*ENCNodegroup, error) {
	if _, ok := enc.Nodegroups[nodegroup]; !ok {
		return &ENCNodegroup{}, errors.New("Nodegroup does not exist")
	}
	parent_chain := node_name + "-" + strings.Join(enc.getParentChain(nodegroup), "-")
	if len(parent_chain) > len(enc.getLongestChain(node_name)) {
		enc.Nodes.Add(parent_chain, parent_chain)
	}

	nodegroup_obj, _ := enc.GetNodegroup(nodegroup)
	nodegroup_obj.Nodes = append(nodegroup_obj.Nodes, node_name)
	enc.Nodegroups[nodegroup] = *nodegroup_obj
	return nodegroup_obj, nil
}

func (enc *ENC) AddNodes(nodegroup string, nodes []string) (*ENCNodegroup, error) {
	for _, node := range nodes {
		_, err := enc.AddNode(nodegroup, node)
		if err != nil {
			return &ENCNodegroup{}, err
		}
	}

	return enc.GetNodegroup(nodegroup)
}

func (enc *ENC) getParentChain(nodegroup_name string) []string {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		panic(err)
	}

	parents := []string{nodegroup_name}
	parent := nodegroup.Parent
	for parent != "" {
		parents = append(parents, parent)
		parent_ng, _ := enc.GetNodegroup(parent)
		parent = parent_ng.Parent
	}

	return parents
}

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

func (enc *ENC) RemoveNode(nodegroup string, node_name string) (*ENCNodegroup, error) {
	if _, ok := enc.Nodegroups[nodegroup]; !ok {
		return &ENCNodegroup{}, errors.New("Nodegroup does not exist")
	}

	chains := enc.Nodes.PrefixSearch(node_name)
	for _, chain := range chains {

		if strings.HasPrefix(chain, node_name+"-"+nodegroup) {
			enc.Nodes.Remove(chain)
			break
		}
	}

	nodegroup_obj, _ := enc.GetNodegroup(nodegroup)
	nodegroup_obj.Nodes = removeByValueSS(nodegroup_obj.Nodes, node_name)
	return nodegroup_obj, nil
}

func (enc *ENC) GetNode(node_name string) (*ENCNodegroup, error) {
	chain := strings.Replace(enc.getLongestChain(node_name), node_name+"-", "", -1)
	path := strings.Split(chain, "-")

	master_nodegroup := &ENCNodegroup{}

	for _, piece := range reverse(path) {
		piece_nodegroup, _ := enc.GetNodegroup(piece)
		master_nodegroup = enc.mergeNodegroups(master_nodegroup, piece_nodegroup)
	}

	return master_nodegroup, nil
}

func (enc *ENC) mergeNodegroups(ng_a *ENCNodegroup, ng_b *ENCNodegroup) *ENCNodegroup {
	new_ng := ENCNodegroup{
		Parent:      ng_b.Parent,
		Nodes:       ng_b.Nodes,
		Environment: ng_a.Environment,
	}

	fmt.Printf("ENV: %#v\n", ng_b.Environment != "")

	if ng_b.Environment != "" {
		new_ng.Environment = ng_b.Environment
	}

	new_ng.Classes = enc.mergeNestedMaps(ng_a.Classes, ng_b.Classes)
	new_ng.Parameters = enc.mergeNestedMaps(ng_a.Parameters, ng_b.Parameters)

	return &new_ng
}

func (enc *ENC) mergeNestedMaps(map_a map[string]interface{}, map_b map[string]interface{}) map[string]interface{} {
	for key, val := range map_b {
		if reflect.ValueOf(val).Kind() == reflect.Map {
			if a_val, ok := map_a[key]; ok && reflect.ValueOf(a_val).Kind() == reflect.Map {
				map_a[key] = enc.mergeNestedMaps(a_val.(map[string]interface{}), val.(map[string]interface{}))
			} else {
				if map_a == nil {
					map_a = make(map[string]interface{})
				}
				map_a[key] = val
			}
		} else {
			if map_a == nil {
				map_a = make(map[string]interface{})
			}
			map_a[key] = val
		}
	}

	return map_a
}

func (enc *ENC) AddParameter(nodegroup_name string, key string, val interface{}) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	nodegroup.Parameters[key] = val
	enc.Nodegroups[nodegroup_name] = *nodegroup
	return nodegroup, nil
}

func (enc *ENC) SetParameter(nodegroup_name string, key string, val interface{}) (*ENCNodegroup, error) {
	return enc.AddParameter(nodegroup_name, key, val)
}

func (enc *ENC) RemoveParameter(nodegroup_name string, key string) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	if _, ok := nodegroup.Parameters[key]; !ok {
		return &ENCNodegroup{}, errors.New("That parameter does not exist for this nodegroup")
	}

	delete(nodegroup.Parameters, key)
	return nodegroup, nil
}

func (enc *ENC) AddClass(nodegroup_name string, key string) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	nodegroup.Classes[key] = make(map[string]interface{})
	enc.Nodegroups[nodegroup_name] = *nodegroup
	return nodegroup, nil
}

func (enc *ENC) SetClasses(nodegroup_name string, key string) (*ENCNodegroup, error) {
	return enc.AddClass(nodegroup_name, key)
}

func (enc *ENC) RemoveClass(nodegroup_name string, key string) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	if _, ok := nodegroup.Classes[key]; !ok {
		return &ENCNodegroup{}, errors.New("That parameter does not exist for this nodegroup")
	}

	delete(nodegroup.Classes, key)
	return nodegroup, nil
}

func (enc *ENC) AddClassParameter(nodegroup_name string, class string, key string, val interface{}) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	if _, ok := nodegroup.Classes[class]; !ok {
		return &ENCNodegroup{}, errors.New("Class does not exist on that nodegroup")
	}

	ng_class := nodegroup.Classes[class].(map[string]interface{})
	ng_class[key] = val
	nodegroup.Classes[class] = ng_class

	return nodegroup, nil
}

func (enc *ENC) RemoveClassParameter(nodegroup_name string, class string, key string) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	if _, ok := nodegroup.Classes[class]; !ok {
		return &ENCNodegroup{}, errors.New("Class does not exist on that nodegroup")
	}

	ng_class := nodegroup.Classes[class].(map[string]interface{})
	if _, ok := ng_class[key]; !ok {
		return &ENCNodegroup{}, errors.New("Parameter for that class does not exist on that nodegroup")
	}

	delete(ng_class, key)
	nodegroup.Classes[class] = ng_class

	return nodegroup, nil
}

func (enc *ENC) SetClassParameter(nodegroup_name string, class string, key string, val interface{}) (*ENCNodegroup, error) {
	return enc.AddClassParameter(nodegroup_name, class, key, val)
}

func (enc *ENC) SetParent(nodegroup_name string, parent string) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	_, parent_err := enc.GetNodegroup(parent)
	if parent_err != nil {
		return &ENCNodegroup{}, errors.New("The parent nodegroup does not exist")
	}

	nodegroup.Parent = parent
	enc.Nodegroups[nodegroup_name] = *nodegroup
	return nodegroup, nil
}

func (enc *ENC) SetEnvironment(nodegroup_name string, env string) (*ENCNodegroup, error) {
	nodegroup, err := enc.GetNodegroup(nodegroup_name)
	if err != nil {
		return &ENCNodegroup{}, err
	}

	nodegroup.Environment = env
	enc.Nodegroups[nodegroup_name] = *nodegroup
	return nodegroup, nil
}
