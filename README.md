# jsonpath

CLI using JSONPath with YAML and JSON. Can be used to find and print subsets of a YAML or JSON file based on a JSONPath
argument. JSON input is handled by first converting to YAML. JSONPath is applied to the YAML and matches printed to
stdout.  Ouput can be YAML or JSON.

JSONPath is not terribly rigorous as standards go. The package used in this project is
[yaml-jsonpath](https://github.com/vmware-labs/yaml-jsonpathhttps://github.com/vmware-labs/yaml-jsonpath). That
project's README file outlines the JSONPath that the project has implemented.

JSONPath is a proposed [IETF standard](https://tools.ietf.org/id/draft-goessner-dispatch-jsonpath-00.html)

| JSONPath         | Description                                                           |
| :--------------- | :-------------------------------------------------------------------- |
| $                | the root object/element                                               |
| @                | the current object/element                                            |
| . or []          | child operator                                                        |
| ..               | nested descendants. JSONPath borrows this syntax from E4X.            |
| *                | wildcard. All objects/elements regardless of their names.             |
| []               | subscript operator. XPath uses it to iterate over element collections |
|                  | and for predicates. In JavaScript and JSON it is the native array     |
|                  | operator.                                                             |
| [,]              | JSONPath allows alternate names or array indices as a set.            |
| [start:end:step] | array slice operator borrowed from ES4.                               |
| ?()              | applies a filter (script) expression.                                 |


| JSONPath                          | Result                                                       |
| :-------------------------------- | :----------------------------------------------------------- |
| $.store.book[*].author            | the authors of all books in the store                        |
| $..author                         | all authors                                                  |
| $.store.*                         | all things in store, which are some books and a red bicycle. |
| $.store..price                    | the price of everything in the store.                        |
| $..book[2]                        | the third book                                               |
| $..book[(@.pages < 10)]           | book with pages value < 10                                   |
| the last book in order.           |                                                              |
| $..book[0,1]                      |                                                              |
| the first two books               |                                                              |
| $..book[?(@.isbn)]                | filter all books with isbn number                            |
| $..book[?(@.price<10)]            | filter all books cheapier than 10                            |
| $..book[(?(@.author=~/Woodward/)] | book with author matching regex                              |
| $..* or $                         | all members of JSON data item.                               |
| $.property~                       | get name of property                                         |

Usage of this tool with a JSONPath will result in zero or more matches. The matches are returned as an array. If there
is an array of length 1 it is printed out on its own. If there are more than one match they are printed out as a YAML or
JSON array. For example, here is an array printed out in JSON format

`$ jsonpath '$..color' -file test/colours.json`
```json
[
 "red",
 "green",
 "blue",
 "cyan",
 "magenta",
 "yellow",
 "black"
]
```

Here is the same query with YAML output

`$ jsonpath '$..color' -file test/colours.json -yaml`
```yaml
- red
- green
- blue
- cyan
- magenta
- yellow
- black
```

This tool uses the [posener completion library](https://github.com/posener/complete/tree/master). You can set it up by
typing `COMP_INSTALL=1 jsonpath`.

I will add a suite of tests to ensure that the library used works and to illustrate as much as possible of jsonpath.

A YAML array can be turned into a list that can be read by bash as an array.

```sh
$ jsonpath -yaml '$..color' -file test/colours.json | awk 'BEGIN {IFS=/\s+/} {printf "%s ", $2}'|xargs
red green blue cyan magenta yellow black
```