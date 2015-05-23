package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hiroosak/gendoc/schema"
)

func validSchemaTree(src string) error {
	if err := isDir(src); err != nil {
		return fmt.Errorf("src is not directory")
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		s, err := schema.NewSchemaFromFile(path, info)
		if err != nil {
			return err
		}

		if js, err := s.ToJSON(); err != nil {
			return nil
		} else {
			return schema.ValidSchema(js)
		}
	})
}
