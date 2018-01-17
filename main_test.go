package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	return f.mode
}

func (f fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f fileInfo) IsDir() bool {
	return f.isDir
}

func (f fileInfo) Sys() interface{} {
	return f.sys
}

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		info     fileInfo
		isGoFile bool
	}{
		{fileInfo{
			"\\foo\\bar\\", 1, os.ModePerm, time.Now(), true, nil,
		}, false},
		{fileInfo{
			"/foo/bar/", 1, os.ModePerm, time.Now(), true, nil,
		}, false},
		{fileInfo{
			"\\foo\\bar\\test.go", 1, os.ModePerm, time.Now(), false, nil,
		}, true},
		{fileInfo{
			"/foo/bar/test.go", 1, os.ModePerm, time.Now(), false, nil,
		}, true},
		{fileInfo{
			"/foo/bar/test.go", 1, os.ModePerm, time.Now(), false, nil,
		}, true},
	}

	for _, test := range tests {
		actual := isGoFile(test.info)
		assert.Equal(t, actual, test.isGoFile)
	}
}
