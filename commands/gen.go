package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateJSON(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("src must be directory")
	}
	if err := createIfNotExist(dst); err != nil {
		return err
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if err := yaml2JSON(path, dst, info); err != nil {
			return err
		}
		return nil
	})
}
