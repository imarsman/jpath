package path

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

// ExpandInterfaceToMatch take interface and expand to match JSON or YAML interface structures
func expandInterfaceToMatch(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = expandInterfaceToMatch(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = expandInterfaceToMatch(v)
		}
	}
	return i
}

// Path a path for a document and associated properties
type Path struct {
	Path     string
	Contents string
	Obj      []interface{}
}

// NewPath create a new path and initialize it
func NewPath(path string, contents string) (p *Path, err error) {
	p = new(Path)
	p.Contents = contents
	p.Path = path

	if IsJSON(p.Contents) {
		p.Contents, err = toYAML(p.Contents)
		if err != nil {
			return
		}
	}

	err = p.process()
	if err != nil {
		return
	}

	return
}

// SubNodeCount the count of sub-documents found
func (p *Path) SubNodeCount() (count int) {
	count = len(p.Obj)

	return
}

// process process a document using the jsonpath
func (p *Path) process() (err error) {
	var parts []string
	parts, err = subNodesStrings(p.Path, p.Contents)
	if err != nil {
		return
	}
	p.Obj, err = fromStringArr(parts)
	if err != nil {
		return
	}

	return
}

// YAML get YAML subset based on jsonpath
func (p *Path) YAML() (y string, err error) {
	if len(p.Obj) == 1 {
		y, err = objToYAML(p.Obj[0])
		if err != nil {
			return
		}
	} else {
		y, err = objToYAML(p.Obj)
		if err != nil {
			return
		}
	}
	y = strings.TrimSpace(y)

	return
}

// JSON get JSON subset based on jsonpath
func (p *Path) JSON() (j string, err error) {
	if len(p.Obj) == 1 {
		j, err = objToJSON(p.Obj[0])
		if err != nil {
			return
		}
	} else {
		j, err = objToJSON(p.Obj)
		if err != nil {
			return
		}
	}
	j = strings.TrimSpace(j)

	return
}

// expandToMatch take interface and recursively update to reflect underlying structure
func expandToMatch(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = expandToMatch(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = expandToMatch(v)
		}
	}
	return i
}

func toYAML(contents string) (output string, err error) {
	if IsJSON(contents) {
		var obj interface{}
		err = json.Unmarshal([]byte(contents), &obj)
		if err != nil {
			return
		}
		obj = expandInterfaceToMatch(obj)
		var bytes []byte
		bytes, err = yaml.Marshal(&obj)
		if err != nil {
			return
		}
		contents = string(bytes)
		return
	}

	output = contents
	return
}

// subNodesStrings get parts of document that match path
func subNodesStrings(path string, content string) (parts []string, err error) {
	var node yaml.Node
	node, err = toYAMLNode(content)
	if err != nil {
		return
	}
	yamlPath, err := yamlpath.NewPath(path)
	if err != nil {
		return
	}

	actual, err := yamlPath.Find(&node)
	if err != nil {
		return
	}
	for _, a := range actual {
		var buf bytes.Buffer
		e := yaml.NewEncoder(&buf)
		e.SetIndent(2)

		err = e.Encode(a)
		if err != nil {
			return
		}
		e.Close()
		parts = append(parts, strings.TrimSpace(buf.String()))
	}

	return
}

// IsJSON test if string is JSON (not exact)
func IsJSON(content string) (isJSON bool) {
	if strings.HasPrefix(content, "{") || strings.HasPrefix(content, "[") {
		if json.Valid([]byte(content)) {
			return true
		}
		fmt.Fprintf(os.Stderr, "Validity check failed")
		os.Exit(1)
	}

	return
}

// toYAMLNode convert string to yaml node
func toYAMLNode(content string) (node yaml.Node, err error) {
	if IsJSON(content) {
		content, err = contentToYAML(content)
		if err != nil {
			return
		}
	}

	err = yaml.Unmarshal([]byte(content), &node)
	if err != nil {
		// fmt.Println("error", err)
		return
	}

	return
}

// toObj convert incoming yaml or json to an interface matching the document
func toObj(content string) (obj interface{}, err error) {
	if IsJSON(content) {
		decoder := json.NewDecoder(strings.NewReader(content))
		err = decoder.Decode(&obj)
		if err != nil {
			return
		}
		obj = expandToMatch(obj)
	} else {
		decoder := yaml.NewDecoder(strings.NewReader(content))
		err = decoder.Decode(&obj)
		if err != nil {
			return
		}
		obj = expandToMatch(obj)
	}

	return
}

// fromStringArr convert an array of sub-document parts to an interface
func fromStringArr(parts []string) (obj []interface{}, err error) {
	for _, part := range parts {
		var o interface{}
		o, err = toObj(part)
		if err != nil {
			return
		}
		obj = append(obj, o)
	}

	return
}

// objToYAML convert an interface to YAML
func objToYAML(obj interface{}) (j string, err error) {
	var bytes []byte
	bytes, err = yaml.Marshal(obj)
	if err != nil {
		return
	}
	j = string(bytes)

	return
}

// objToJSON convert an interface to JSON
func objToJSON(obj interface{}) (j string, err error) {
	var bytes []byte
	bytes, err = json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return
	}
	j = string(bytes)

	return
}

// contentToYAML convert content to YAML
func contentToYAML(content string) (y string, err error) {
	str := string(content)
	str = strings.TrimSpace(str)

	var obj interface{}

	obj, err = toObj(content)
	if err != nil {
		return
	}

	var bytes []byte
	bytes, err = yaml.Marshal(obj)
	y = string(bytes)

	return
}

// contentToYAML convert content to YAML
func contentToJSON(content string) (j string, err error) {
	str := string(content)
	str = strings.TrimSpace(str)

	var obj interface{}

	obj, err = toObj(content)
	if err != nil {
		return
	}

	var bytes []byte
	bytes, err = json.Marshal(obj)
	j = string(bytes)

	return
}
