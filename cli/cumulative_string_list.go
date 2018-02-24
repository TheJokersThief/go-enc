package cli

import (
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

type stringList []string

func (i *stringList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *stringList) String() string {
	return strings.Join(*i, ", ")
}

func (i *stringList) IsCumulative() bool {
	return true
}

func StringList(s kingpin.Settings) (target *[]string) {
	target = new([]string)
	s.SetValue((*stringList)(target))
	return
}
