package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// main options
	config = flag.String("c", "", "config file")

	layers = map[string]int{}
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go-clean-arch [flags] [path]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if err := readConfig(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "error: no arguments specified, expecting one")
		return
	}

	path, _ := filepath.Abs(flag.Arg(0))
	switch dir, err := os.Stat(path); {
	case err != nil:
		fmt.Fprintf(os.Stderr, err.Error())
		return
	case dir.IsDir():
		processor := NewProcessor(path)
		processor.Process()
	default:
		fmt.Fprintf(os.Stderr, "error: can not use go-clean-arch on a single file")
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

	file, err := os.Open(*config)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &layers)
}
