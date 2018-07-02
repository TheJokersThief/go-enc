package enc

import (
	"fmt"
	"reflect"
)

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func removeByValueSS(a []string, val string) []string {
	newArray := make([]string, 0)
	for _, x := range a {
		if x != val {
			newArray = append(newArray, x)
		}
	}

	return newArray
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func stringifyYAMLMapKeys(in interface{}) interface{} {
	switch in := in.(type) {
	case []interface{}:
		res := make([]interface{}, len(in))
		for i, v := range in {
			res[i] = stringifyYAMLMapKeys(v)
		}
		return res
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		for k, v := range in {
			res[fmt.Sprintf("%v", k)] = stringifyYAMLMapKeys(v)
		}
		return res
	default:
		return in
	}
}

// mergeNestedMaps travels nested maps and adds the values from mapB into mapA (overwriting
// what exists and preserving what's unique in both)
func MergeNestedMaps(mapA map[string]interface{}, mapB map[string]interface{}) map[string]interface{} {
	for key, val := range mapB {
		if reflect.ValueOf(val).Kind() == reflect.Map {
			if aVal, ok := mapA[key]; ok && reflect.ValueOf(aVal).Kind() == reflect.Map {
				mapA[key] = MergeNestedMaps(aVal.(map[string]interface{}), val.(map[string]interface{}))
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
