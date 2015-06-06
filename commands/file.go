package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/hiroosak/gendoc/schema"
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

// yaml2JSON creates json schema file from yaml file.
func yaml2JSON(src, dst string, info os.FileInfo) error {
	ext := filepath.Ext(src)
	if ext != ".yaml" && ext != ".yml" {
		return nil
	}
	r, err := schema.NewSchemaFromFile(src, info)
	if err != nil {
		return err
	}
	j, err := r.ToJSON()
	if err != nil {
		return err
	}

	base := filepath.Base(src)
	dstFile := base[0:len(base)-len(ext)] + ".json"
	dstPath := path.Join(dst, dstFile)

	ioutil.WriteFile(dstPath, j, filePerm)

	return nil
}

// isDir returns true if path is a directory.
func isDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not such directory")
	}
	return nil
}

// createIfNotExist creates the directory to path if it doesn't exist.
func createIfNotExist(path string) error {
	err := isDir(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, dirPerm)
		return nil
	}
	return err
}
