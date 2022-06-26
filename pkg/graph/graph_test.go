package graph

import (
	"testing"

	"gopkg.in/yaml.v3"
)

var colours = `[
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
]`

func TestYaml(t *testing.T) {
	var yamlNode yaml.Node
	err := yaml.Unmarshal([]byte(colours), &yamlNode)
	if err != nil {
		return
	}
	kind := yamlNode.Kind

	t.Log(kind == yaml.DocumentNode)
}
