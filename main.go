package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// main options
	config = flag.String("c", "", "config file")
)

var (
	fileSet  = token.NewFileSet()
	layers   = map[string]int{}
	basePath = ""
)

// isGoFile determines whether a file is a Go file.
func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}
	// Don't complain if a file was deleted in the meantime (i.e.
	// the directory changed concurrently while running gofmt).
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Encountered error: %v\n", err)
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go-clean-arch [flags] [path ...]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if err := readConfig(); err != nil {
		scanner.PrintError(os.Stderr, err)
		return
	}

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "error: no arguments specified.")
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		basePath, _ = filepath.Abs(path)
		switch dir, err := os.Stat(path); {
		case err != nil:
			scanner.PrintError(os.Stderr, err)
			return
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				scanner.PrintError(os.Stderr, err)
				return
			}
		}
	}
	os.Exit(0)
}

func readConfig() error {

	if *config == "" {
		return errors.New("error: config file was not specified")
	}

	if _, err := os.Stat(*config); err != nil {
		return err
	}

	file, err := ioutil.ReadFile(*config)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &layers)
}

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {

	filename, _ = filepath.Abs(filename)
	packagePath := getPackage(filename)
	cleanArchLayerIndex := getCleanArchLayerIndex(packagePath)
	if cleanArchLayerIndex == 0 {
		return nil
	}

	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	file, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		importLayerIndex := getCleanArchLayerIndex(importPath)
		if importLayerIndex > cleanArchLayerIndex {
			fmt.Printf("error: bad dependency on '%s' in layer '%s' ('%s')\n", importPath, packagePath, filename)
		}
	}

	return nil
}

func getCleanArchLayerIndex(importPath string) int {
	for k, v := range layers {
		length := len(k)
		if len(importPath) >= length && importPath[0:length] == k {
			return v
		}
	}
	return 0
}

func getPackage(filename string) string {
	relativePath := strings.Replace(filename, basePath, "", 1)
	relativePath = strings.TrimLeft(relativePath, string(os.PathSeparator))
	relativePath = filepath.Dir(relativePath)
	relativePath = filepath.ToSlash(relativePath)
	return relativePath
}
