package enc

import ()

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func removeByValueSS(a []string, val string) []string {
	new_array := make([]string, 0)
	for _, x := range a {
		if x != val {
			new_array = append(new_array, x)
		}
	}

	return new_array
}