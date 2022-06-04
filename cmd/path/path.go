package path

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

type document struct {
}

// SubNodesStrings get parts of document that match path
func SubNodesStrings(path string, content string) (parts []string, err error) {
	var node yaml.Node
	node, err = ToYAMLNode(content)
	if err != nil {
		return
	}
	p, err := yamlpath.NewPath(path)
	if err != nil {
		return
	}

	actual, err := p.Find(&node)
	for _, a := range actual {
		var buf bytes.Buffer
		e := yaml.NewEncoder(&buf)
		e.SetIndent(2)

		err = e.Encode(a)
		e.Close()
		parts = append(parts, strings.TrimSpace(buf.String()))
	}

	return
}

func isJSON(content string) bool {
	return strings.HasPrefix(content, "{") || strings.HasPrefix(content, "[")
}

func ToYAMLNode(content string) (node yaml.Node, err error) {
	if isJSON(content) {
		content, err = ToYAML(content)
		if err != nil {
			return
		}
	}

	err = yaml.Unmarshal([]byte(content), &node)
	if err != nil {
		return
	}

	return
}

func ToYAML(content string) (y string, err error) {
	str := string(content)
	str = strings.TrimSpace(str)

	if isJSON(content) {
		obj := make(map[string]interface{})
		fmt.Println(obj)
		e := json.NewDecoder(strings.NewReader(string(content)))
		err = e.Decode(&obj)
		if err != nil {
			return
		}
		var bytes []byte
		bytes, err = yaml.Marshal(obj)
		if err != nil {
			return
		}
		y = string(bytes)
		return
	}
	y = string(content)

	return
}
