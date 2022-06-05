package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/imarsman/jsonpath/cmd/path"
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

type Args struct {
	JSON bool   `arg:"-j,--json" help:"output json"`
	YAML bool   `arg:"-y,--yaml" help:"output yaml"`
	Path string `arg:"required,positional" help:"jsonpath to use"`
	File string `arg:"-f,--file" help:"file to use instead of stdin"`
}

// Version get version information
func (Args) Version() string {
	var buf = new(bytes.Buffer)

	msg := "jsonpath"
	buf.WriteString(fmt.Sprintln(msg))
	buf.WriteString(fmt.Sprintln(strings.Repeat("-", len(msg))))

	if GitCommit != "" {
		buf.WriteString(fmt.Sprintf("Commit: %13s\n", GitCommit))
	}
	if Date != "" {
		buf.WriteString(fmt.Sprintf("Compile Date: %12s\n", Date))
	}

	return buf.String()
}

func main() {
	var args = Args{}
	// Make config to hold various parameters
	cmd := &complete.Command{
		Flags: map[string]complete.Predictor{
			// "path": predict.Nothing,
			"yaml": predict.Nothing,
			"json": predict.Nothing,
			"file": predict.Files("./*"),
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
	path, err := path.NewPath(args.Path, content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if args.YAML {
		contents, err := path.YAML()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(contents)
	} else if args.JSON {
		contents, err := path.JSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(contents)
	} else {
		contents, err := path.JSON()
		fmt.Println("here", err)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(contents)
	}
}
