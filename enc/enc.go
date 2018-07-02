package enc

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/derekparker/trie"
)

var (
	CHAIN_SEPARATION_CHARACTER = "$$"
)

// Nodegroup represents groups of nodes and meta information about them
type Nodegroup struct {
	Parent      string                 `json:"parent,omitempty" yaml:"parent,omitempty"`
	Classes     map[string]interface{} `json:"classes,omitempty" yaml:"classes,omitempty"`
	Nodes       []string               `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Environment string                 `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// ENC represents the entire structure of the External Node Classifier
type ENC struct {
	Name       string
	Nodegroups map[string]Nodegroup
	Nodes      *trie.Trie
	ConfigType string
	FileName   string
	ConfigLink *Config
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

	return &Nodegroup{}, fmt.Errorf("Nodegroup does not exist: %s", name)
}

// GetNodegroup retrieves a nodegroup by name
func (enc *ENC) GetNodegroup(nodegroupName string) (*Nodegroup, error) {
	var (
		nodegroup         string
		cluster           string
		clusterNodegroups map[string]Nodegroup
	)

	config := enc.ConfigLink

	if config != nil {
		if strings.Contains(nodegroupName, "@") {
			nodegroupSplit := strings.Split(nodegroupName, "@")
			nodegroup, cluster = nodegroupSplit[0], nodegroupSplit[1]
		} else {
			nodegroup, cluster = nodegroupName, enc.Name
		}

		clusterNodegroups = config.ENCs[cluster].Nodegroups

		keys := make([]string, 0)
		for key, _ := range clusterNodegroups {
			keys = append(keys, key)
		}
	} else {
		clusterNodegroups, nodegroup = enc.Nodegroups, nodegroupName
	}

	if val, ok := clusterNodegroups[nodegroup]; ok {
		return &val, nil
	}

	return &Nodegroup{}, fmt.Errorf("Nodegroup does not exist: %s", nodegroupName)
}

// AddNode adds a single node to a nodegroup
func (enc *ENC) AddNode(nodegroup string, nodeName string) (*Nodegroup, error) {
	if _, ok := enc.Nodegroups[nodegroup]; !ok {
		return &Nodegroup{}, errors.New("Nodegroup does not exist")
	}
	parentChain := nodeName + CHAIN_SEPARATION_CHARACTER + strings.Join(reverse(enc.getParentChain(nodegroup)), CHAIN_SEPARATION_CHARACTER)

	if _, ok := enc.Nodes.Find(nodeName); !ok {
		enc.Nodes.Add(nodeName, nodeName)
	}
	enc.Nodes.Add(parentChain, parentChain)

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

func (enc *ENC) ParentChainWrapper(nodegroupName string) []string {
	return enc.getParentChain(nodegroupName)
}

// getParentChain generates a path of parents until it reaches the top
func (enc *ENC) getParentChain(nodegroupName string) []string {
	nodegroup, err := enc.GetNodegroup(nodegroupName)
	errCheck(err)

	// If nodegroup doesn't have an explicit cluster, it's the current cluster
	if !strings.Contains(nodegroupName, "@") && enc.Name != "" {
		nodegroupName = nodegroupName + "@" + enc.Name
	}

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

	// If nodegroup doesn't have an explicit cluster, it's the current cluster
	if !strings.Contains(nodegroup, "@") && enc.Name != "" {
		nodegroup = nodegroup + "@" + enc.Name
	}

	children := enc.Nodes.FuzzySearch(nodegroup + CHAIN_SEPARATION_CHARACTER)
	if len(children) == 0 {
		chains := enc.Nodes.PrefixSearch(nodeName)
		for _, chain := range chains {

			if strings.HasSuffix(chain, nodegroup) {
				enc.Nodes.Remove(chain)
				break
			}
		}
	}

	nodegroupObj, _ := enc.GetNodegroup(nodegroup)
	nodegroupObj.Nodes = removeByValueSS(nodegroupObj.Nodes, nodeName)
	return nodegroupObj, nil
}

// GetNode retrieves a nodegroup that represents all inherited values for a node
func (enc *ENC) GetNode(nodeName string) (*Nodegroup, error) {
	var (
		matchedNodegroups []*Nodegroup
	)

	chains, err := enc.GetChains(nodeName)
	errCheck(err)

	commonChain, alteredChains := enc.findCommonChain(chains)
	masterNodegroup := &Nodegroup{}

	// Find merges for all the leafs of the trie
	for _, chain := range alteredChains {
		if chain != "" {
			matchedNodegroups = append(matchedNodegroups, enc.getMergedChainNodegroup(chain))
		}
	}

	if len(matchedNodegroups) > 1 {
		masterNodegroup, err = enc.ConflictMerge(matchedNodegroups)
		errCheck(err)
	}

	// Finally, get the info for the common chain and merge the final data onto it
	masterNodegroup = enc.mergeNodegroups(enc.getMergedChainNodegroup(commonChain), masterNodegroup)

	return masterNodegroup, nil
}

func (enc *ENC) getMergedChainNodegroup(chain string) *Nodegroup {
	path := strings.Split(chain, CHAIN_SEPARATION_CHARACTER)

	masterNodegroup := &Nodegroup{}

	for _, piece := range path {
		pieceNodegroup, _ := enc.GetNodegroup(piece)
		masterNodegroup = enc.mergeNodegroups(masterNodegroup, pieceNodegroup)
	}

	return masterNodegroup
}

// Returns the common chain (in ALL chains) and the chains stripped of the common chain
func (enc *ENC) findCommonChain(chains []string) (string, []string) {
	var (
		commonPieces  []string
		alteredChains []string
	)

	firstChain := chains[0]
	pieces := strings.Split(firstChain, CHAIN_SEPARATION_CHARACTER)

	for _, piece := range pieces {
		isCommon := true

		currentPiece := strings.Join(append(commonPieces, piece), CHAIN_SEPARATION_CHARACTER)
		for _, chain := range chains {
			if !strings.HasPrefix(chain, currentPiece) {
				isCommon = false
			}
		}

		if isCommon {
			commonPieces = append(commonPieces, piece)
		} else {
			// If it's not common, stop looking any further
			break
		}
	}

	commonChain := strings.Join(commonPieces, CHAIN_SEPARATION_CHARACTER)

	for _, chain := range chains {
		// Remove with separation character first
		chain = strings.Replace(chain, commonChain+CHAIN_SEPARATION_CHARACTER, "", -1)
		// Remove without separation character second
		chain = strings.Replace(chain, commonChain, "", -1)
		alteredChains = append(alteredChains, chain)
	}

	return commonChain, alteredChains
}

func (enc *ENC) ConflictMerge(nodegroups []*Nodegroup) (*Nodegroup, error) {
	xNG := &Nodegroup{}

	for _, yNG := range nodegroups {
		// If both ngs have the same class, if they have the same param,
		// the values must be the same
		for yClass, yParams := range yNG.Classes {
			if xClass, hasClass := xNG.Classes[yClass]; hasClass {
				if yParamsConverted, convertOk := yParams.(map[string]interface{}); convertOk {
					for yKey, yVal := range yParamsConverted {
						if xVal, hasKey := xClass.(map[string]interface{})[yKey]; hasKey {
							if !reflect.DeepEqual(xVal, yVal) {
								return &Nodegroup{}, fmt.Errorf(
									"Conflict detected: [class %s, key %s, xVal: %#v, yVal: %#v]",
									yClass, yKey, xVal, yVal)
							}
						}
					}
				}
			}
		}

		// If both ngs have the same parameter, if they have the same option,
		// the values must be the same
		for yParameter, yOptions := range yNG.Parameters {
			if xParameter, hasParameter := xNG.Parameters[yParameter]; hasParameter {
				if yOptionsConverted, convertOk := yOptions.(map[string]interface{}); convertOk {
					for yKey, yVal := range yOptionsConverted {
						if xVal, hasKey := xParameter.(map[string]interface{})[yKey]; hasKey {
							if !reflect.DeepEqual(xVal, yVal) {
								return &Nodegroup{}, fmt.Errorf(
									"Conflict detected: [parameter %s, key %s, xVal: %#v, yVal: %#v]",
									yParameter, yKey, xVal, yVal)
							}
						}
					}
				}
			}
		}

		// Shouldn't matter which way these are merged
		xNG = enc.mergeNodegroups(xNG, yNG)
	}

	return xNG, nil
}

// Get all possible parents for a node from the trie
func (enc *ENC) GetChains(nodeName string) ([]string, error) {
	root, ok := enc.Nodes.Find(nodeName)
	if !ok {
		return []string{}, fmt.Errorf("Could not find node in ENC")
	}

	return enc.travelChain(root.Parent(), nodeName), nil
}

func (enc *ENC) travelChain(root *trie.Node, currentChain string) []string {
	var trackerChain []string
	for letter, node := range root.Children() {
		if len(root.Children()) > 1 && string(letter) == "\x00" {
			continue
		}

		newChain := currentChain
		if string(letter) != "\x00" {
			newChain = newChain + string(letter)
		}

		if len(node.Children()) > 0 {
			childChains := enc.travelChain(node, newChain)
			trackerChain = append(trackerChain, childChains...)
		} else {
			trackerChain = append(trackerChain, newChain)
		}
	}

	return trackerChain
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

	newNG.Classes = MergeNestedMaps(ngA.Classes, ngB.Classes)
	newNG.Parameters = MergeNestedMaps(ngA.Parameters, ngB.Parameters)

	return &newNG
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
