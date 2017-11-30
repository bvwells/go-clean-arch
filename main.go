package main

import (
	"bufio"
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
	f, err := os.Open(*config)
	if err != nil {
		return err
	}
	defer f.Close()
	
	layerIndex := 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		layers[scanner.Text()] = layerIndex
		layerIndex++
	}
	return scanner.Err()
}

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {

	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}
	filename, _ = filepath.Abs(filename)
	cleanArchLayer, err := getBaseDirectory(filename)
	if err != nil {
		return nil
	}
	cleanArchLayerIndex, okay := layers[cleanArchLayer]
	if !okay {
		return nil
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	//	fmt.Printf("Parsing file: %s\n", filename)
	file, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		if importLayer, err := getBasePackage(importPath); err != nil {
			continue
		} else {
			fmt.Printf("comparing clean arch package %s to import %s\n", cleanArchLayer, importLayer)

			importLayerIndex, found := layers[importLayer]
			if !found {
				continue
			} else {
				if importLayerIndex > cleanArchLayerIndex {
					fmt.Printf("Error in clean architecture in file %s.\n!", filename)
				}
			}
		}
	}

	return nil
}

func getBasePackage(importPath string) (string, error) {
	index := strings.Index(importPath, "/")
	if index == -1 {
		return "", errors.New("base directory does not exist")
	}
	return importPath[0:index], nil
}

func getBaseDirectory(filename string) (string, error) {
	relativePath := strings.Replace(filename, basePath, "", 1)
	relativePath = strings.TrimLeft(relativePath, string(os.PathSeparator))
	index := strings.Index(relativePath, string(os.PathSeparator))
	// return if no base directory.
	if index == -1 {
		return "", errors.New("base directory does not exist")
	}
	return relativePath[0:index], nil
}
