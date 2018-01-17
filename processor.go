package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Processor defines a processor for processing.
type Processor struct {
	basePath string
}

// NewProcessor creates a new processor.
func NewProcessor(basePath string) *Processor {
	return &Processor{basePath}
}

// Process processed all the go files relative to the base path to
// determine whether clean architecture has been violated.
func (p *Processor) Process() error {
	return p.walkDir(p.basePath)
}

func (p *Processor) walkDir(path string) error {
	return filepath.Walk(path, p.visitFile)
}

// isGoFile determines whether a file is a Go file.
func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func (p *Processor) visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = p.processFile(path)
	}
	// Don't complain if a file was deleted in the meantime (i.e.
	// the directory changed concurrently while running).
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Encountered error: %v\n", err)
	}
	return nil
}

func (p *Processor) processFile(filename string) error {
	filename, _ = filepath.Abs(filename)
	packagePath := p.getPackage(filename)
	cleanArchLayerIndex := getCleanArchLayerIndex(packagePath)
	if cleanArchLayerIndex == 0 {
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, src, parser.ImportsOnly)
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

func (p *Processor) getPackage(filename string) string {
	relativePath := strings.Replace(filename, p.basePath, "", 1)
	relativePath = strings.TrimLeft(relativePath, string(os.PathSeparator))
	relativePath = filepath.Dir(relativePath)
	relativePath = filepath.ToSlash(relativePath)
	return relativePath
}
