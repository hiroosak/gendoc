package commands

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateIfNotExist(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	target := dir + "/foo"

	if err := createIfNotExist(target); err != nil {
		t.Fatal(err.Error())
	}
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err.Error())
	}
}
