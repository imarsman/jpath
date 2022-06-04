package main

import (
	"bytes"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

func main() {
	var n yaml.Node
	p, err := yamlpath.NewPath("")
	if err != nil {

	}
	actual, err := p.Find(&n)

	actualStrings := []string{}
	for _, a := range actual {
		var buf bytes.Buffer
		e := yaml.NewEncoder(&buf)
		e.SetIndent(2)

		err = e.Encode(a)
		e.Close()
		actualStrings = append(actualStrings, buf.String())
	}
}
