package path

import (
	"testing"

	"github.com/matryer/is"
)

var j = `{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}
`

var j2 = `[
	{
		"color": "red",
		"value": "#f00"
	},
	{
		"color": "green",
		"value": "#0f0"
	},
	{
		"color": "blue",
		"value": "#00f"
	},
	{
		"color": "cyan",
		"value": "#0ff"
	},
	{
		"color": "magenta",
		"value": "#f0f"
	},
	{
		"color": "yellow",
		"value": "#ff0"
	},
	{
		"color": "black",
		"value": "#000"
	}
]
`

var y = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
      - name: nginy
        image: nginy
        ports:
        - containerPort: 81
`

// func TestIsJSON(t *testing.T) {
// 	is := is.New(t)

// 	testIsJSON := isJSON(j)
// 	is.True(testIsJSON)

// 	testIsJSON = isJSON(y)
// 	is.True(testIsJSON == false)
// }

// func TestJSONToYAML(t *testing.T) {
// 	is := is.New(t)

// 	testIsJSON := isJSON(j)
// 	is.True(testIsJSON)

// 	yaml, err := ContentToYAML(j)
// 	is.NoErr(err)
// 	t.Log(yaml)
// }

// func TestJSONPath(t *testing.T) {
// 	is := is.New(t)
// 	path := "$..spec.containers[*].image"

// 	parts, err := SubNodesStrings(path, y)
// 	is.NoErr(err)
// 	t.Log("parts length", len(parts))
// 	// for _, p := range parts {
// 	// 	t.Log(p)
// 	// }
// 	out, err := toYAML(parts)
// 	is.NoErr(err)
// 	t.Log("yaml", out)
// 	out, err = toYAML(parts)
// 	is.NoErr(err)
// 	t.Log("yaml", out)
// }

// func TestJSONPathFromJSON(t *testing.T) {
// 	is := is.New(t)
// 	path := "$.glossary..GlossDiv.title"
// 	parts, err := SubNodesStrings(path, j)
// 	is.NoErr(err)
// 	t.Log("parts length", len(parts))
// 	// out, err := toJSON(parts)
// 	// t.Log("json", out)

// 	path = "$..GlossEntry"
// 	parts, err = SubNodesStrings(path, j)
// 	is.NoErr(err)
// 	t.Log("parts length", len(parts))
// 	// out, err = toJSON(parts)
// 	// t.Log("json", out)
// 	t.Log("parts", parts)

// 	path = "$..GlossEntry..Acronym"
// 	parts, err = SubNodesStrings(path, j)
// 	is.NoErr(err)
// 	t.Log("parts length", len(parts))
// 	// out, err = toJSON(parts)
// 	// t.Log("json", out)
// 	t.Log("parts", parts)
// 	// out, err = toYAML(parts)
// 	// t.Log("yaml", out)
// 	t.Log(parts)
// }

func TestJSONPathFromJSON2(t *testing.T) {
	is := is.New(t)
	path := "$..[?(@.color=~/red/)].color"
	parts, err := subNodesStrings(path, j2)
	is.NoErr(err)
	obj, err := fromStringArr(parts)
	output, err := objToJSON(obj)
	is.NoErr(err)
	t.Log(output)
	output, err = objToYAML(obj)
	is.NoErr(err)
	t.Log(output)
}

func TestJSONPathFromJSON2Longer(t *testing.T) {
	is := is.New(t)
	path := "$..color"
	parts, err := subNodesStrings(path, j2)
	is.NoErr(err)
	obj, err := fromStringArr(parts)
	output, err := objToJSON(obj)
	is.NoErr(err)
	t.Log(output)
	output, err = objToYAML(obj)
	is.NoErr(err)
	t.Log(output)
}

func TestLongFromJSON2(t *testing.T) {
	is := is.New(t)
	path := "$.spec.template"

	parts, err := subNodesStrings(path, y)
	is.NoErr(err)
	obj, err := fromStringArr(parts)
	output, err := objToJSON(obj)
	is.NoErr(err)
	t.Log(output)
	output, err = objToYAML(obj)
	is.NoErr(err)
	t.Log(output)
}
