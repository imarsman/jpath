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

func TestIsJSON(t *testing.T) {
	is := is.New(t)

	testIsJSON := isJSON(j)
	is.True(testIsJSON)

	testIsJSON = isJSON(y)
	is.True(testIsJSON == false)
}

func TestJSONToYAML(t *testing.T) {
	is := is.New(t)

	testIsJSON := isJSON(j)
	is.True(testIsJSON)

	yaml, err := ToYAML(j)
	is.NoErr(err)
	t.Log(yaml)
}

func TestJSONPath(t *testing.T) {
	is := is.New(t)
	path := "$..spec.containers[*].image"

	parts, err := SubNodesStrings(path, y)
	is.NoErr(err)
	t.Log("parts length", len(parts))
	for _, p := range parts {
		t.Log(p)
	}
}
