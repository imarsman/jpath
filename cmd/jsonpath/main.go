package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/imarsman/jsonpath/pkg/path"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

// GitCommit the latest git commit tag
var GitCommit string

// GitLastTag last tag committed
var GitLastTag string

// GitExactTag exact last tag
var GitExactTag string

// Date latest compile date
var Date string

// Args cli args
type Args struct {
	JSON    bool   `arg:"-j,--json" help:"output json"`
	YAML    bool   `arg:"-y,--yaml" help:"output yaml"`
	Path    string `arg:"positional" help:"jsonpath to use"`
	File    string `arg:"-f,--file" help:"file to use instead of stdin"`
	Type    bool   `arg:"-t,--type" help:"Show inferred type of input"`
	Summary bool   `arg:"-s,--summary" help:"Show summary"`
}

// Version get version information
func (Args) Version() string {
	var buf = new(bytes.Buffer)

	msg := "jsonpath"
	buf.WriteString(fmt.Sprintln(msg))
	buf.WriteString(fmt.Sprintln(strings.Repeat("-", len(msg))))

	if GitCommit != "" {
		buf.WriteString(fmt.Sprintf("Commit: %8s\n", GitCommit))
	}
	if Date != "" {
		buf.WriteString(fmt.Sprintf("Date: %23s\n", Date))
	}
	if GitExactTag != "" {
		buf.WriteString(fmt.Sprintf("Tag: %11s\n", GitExactTag))
	}
	buf.WriteString(fmt.Sprintf("OS: %11s\n", runtime.GOOS))
	buf.WriteString(fmt.Sprintf("ARCH: %8s\n", runtime.GOARCH))

	return buf.String()
}

func main() {
	var args = Args{}
	// Make config to hold various parameters
	cmd := &complete.Command{
		Flags: map[string]complete.Predictor{
			// "path": predict.Nothing,
			"yaml":    predict.Nothing,
			"json":    predict.Nothing,
			"file":    predict.Files("./*"),
			"type":    predict.Nothing,
			"summary": predict.Nothing,
		},
	}
	cmd.Complete("jsonpath")

	arg.MustParse(&args)

	var lines []string
	var content string

	if args.File != "" {
		dir, file := filepath.Split(args.File)
		pathToLoad := filepath.Join(dir, file)
		if _, err := os.Stat(pathToLoad); err == nil {
			bytes, err := ioutil.ReadFile(pathToLoad)
			if err != nil {

			}
			content = string(bytes)
		} else if errors.Is(err, os.ErrNotExist) {
			err = fmt.Errorf("file %s not found", args.File)
			fmt.Println(err)
			os.Exit(1)
		} else {
			err = fmt.Errorf("file %s not found", args.File)
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			var scanner = bufio.NewScanner(os.Stdin)

			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				item := scanner.Text()
				item = strings.TrimSpace(item)

				lines = append(lines, item)
			}
			content = strings.Join(lines, "\n")
		}
	}

	// jsonpath library will deal with "" as if it is "$"
	if args.Path == "" {
		args.Path = "$"
	}

	// Print out type based on guess tied to JSON { and [ starting characters
	if args.Type {
		if json.Valid([]byte(content)) {
			fmt.Println("JSON")
		} else {

			fmt.Println("YAML")
		}
		os.Exit(0)
	}

	// Handle jsonpath string
	path, err := path.NewPath(args.Path, content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Print out results as YAML or as JSON
	if args.YAML {
		contents, err := path.YAML()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(contents)
		if args.Summary {
			if path.SubNodeCount() == 1 {
				fmt.Printf("\n%d item found\n", path.SubNodeCount())
			} else {
				fmt.Printf("\n%d items found\n", path.SubNodeCount())
			}
		}
	} else if args.JSON {
		contents, err := path.JSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(contents)
		if path.SubNodeCount() == 1 {
			fmt.Printf("\n%d item found\n", path.SubNodeCount())
		} else {
			fmt.Printf("\n%d items found\n", path.SubNodeCount())
		}
	} else {
		// Print out json by default
		contents, err := path.JSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(contents)
		if path.SubNodeCount() == 1 {
			fmt.Printf("\n%d item found\n", path.SubNodeCount())
		} else {
			fmt.Printf("\n%d items found\n", path.SubNodeCount())
		}
	}
}
