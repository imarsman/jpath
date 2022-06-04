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

func toYAML(parts []string) (y string, err error) {
	bytes, err := yaml.Marshal(&parts)
	if err != nil {
		return
	}
	y = string(bytes)
	fmt.Println(y)

	return
}

func toJSON(parts []string) (j string, err error) {
	bytes, err := json.MarshalIndent(&parts, "", "  ")
	if err != nil {
		return
	}
	j = string(bytes)
	fmt.Println(j)

	return
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

func ToCompact(parts []string) string {
	var output = []string{}

	for _, p := range parts {
		if !strings.Contains(p, "\n") {
			output = append(output, `"`+p+`"`)
		} else {
			p = strings.ReplaceAll(p, "\n", "\\n")
			output = append(output, `"`+p+`"`)
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(output, ", "))
}

func jsonToYaml(content string) (y string, err error) {
	var bytes []byte

	if isJSON(content) {
		obj := make(map[string]interface{})
		var objArr []map[string]interface{}

		e := json.NewDecoder(strings.NewReader(string(content)))
		if strings.HasPrefix(content, "[") {
			err = e.Decode(&objArr)
			if err != nil {
				return
			}
			bytes, err = yaml.Marshal(objArr)
			if err != nil {
				return
			}
			y = string(bytes)

		} else {
			err = e.Decode(&obj)
			if err != nil {
				return
			}
			bytes, err = yaml.Marshal(obj)
			if err != nil {
				return
			}
			y = string(bytes)
		}
	}
	return
}

func ToYAML(content string) (y string, err error) {
	str := string(content)
	str = strings.TrimSpace(str)

	content, err = jsonToYaml(content)
	if err != nil {
		return
	}
	y = string(content)

	return
}
