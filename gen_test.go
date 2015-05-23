package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGenerateJSON(t *testing.T) {
	src, err := ioutil.TempDir("", "src")
	if err != nil {
		t.Fatal(err)
	}
	dst, err := ioutil.TempDir("", "dst")
	if err != nil {
		t.Fatal(err)
	}
	srcfile := path.Join(src, "u.yaml")
	if err := ioutil.WriteFile(srcfile, []byte(schemaTmpl), 0755); err != nil {
		t.Fatal(err)
	}
	if err := generateJSON(src, dst); err != nil {
		t.Error(err)
	}
	dstfile := path.Join(dst, "u.json")
	if _, err = os.Stat(dstfile); err != nil {
		t.Error(err)
	}
	if _, err := ioutil.ReadFile(dstfile); err != nil {
		t.Error(err)
	}
}
