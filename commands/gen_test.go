package commands

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
	srcfile := path.Join(src, "user.yaml")
	w := renderScaffold("user")

	if err := ioutil.WriteFile(srcfile, w.Bytes(), 0755); err != nil {
		t.Fatal(err)
	}
	if err := GenerateJSON(src, dst); err != nil {
		t.Error(err)
	}
	dstfile := path.Join(dst, "user.json")
	if _, err = os.Stat(dstfile); err != nil {
		t.Error(err)
	}
	if _, err := ioutil.ReadFile(dstfile); err != nil {
		t.Error(err)
	}
}
