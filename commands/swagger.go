package commands

import (
	"bytes"
	"os"

	"github.com/hiroosak/gendoc/swagger"
)

func GenerateSwaggerJSON(src, dst string) error {
	if err := isDir(src); err != nil {
		return err
	}
	resources, err := readResources(src)
	if err != nil {
		return err
	}

	c := swagger.NewConverter(swagger.Param{
		Title:       "test",
		Version:     "1.0.0",
		Description: "",
		BaseURL:     "http://localhost/api/v1",
	})

	w := bytes.NewBuffer([]byte{})
	if err := c.Convert(w, resources); err != nil {
		return err
	}
	_, err = w.WriteTo(os.Stdout)
	return err
}
