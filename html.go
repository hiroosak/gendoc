package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hiroosak/gendoc/schema"
)

type htmlParam struct {
	SchemaSlice schema.SchemaSlice
	Meta        Meta
	Overview    template.HTML
}

func generateHTML(src, metafile, overviewfile, templatePath string) error {
	if err := isDir(src); err != nil {
		return err
	}
	resources, err := readResources(src)
	if err != nil {
		return err
	}
	meta, err := readMeta(metafile)
	if err != nil {
		return err
	}
	overview := readOverview(overviewfile)

	param := htmlParam{
		SchemaSlice: resources,
		Meta:        meta,
		Overview:    overview,
	}

	files := templatePath + "/*.tpl"
	tmpl := template.Must(template.New("root").
		Funcs(generateFuncMap(meta)).
		ParseGlob(files))

	w := bytes.NewBuffer([]byte{})
	if err := tmpl.ExecuteTemplate(w, "root", param); err != nil {
		return err
	}

	w.WriteTo(os.Stdout)

	return nil
}

func readResources(src string) (schema.SchemaSlice, error) {
	var resources schema.SchemaSlice
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if r, err := schema.NewSchemaFromFile(path, info); err == nil {
			resources = append(resources, *r)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return resources, err
}

func readOverview(src string) template.HTML {
	var overview template.HTML
	if src != "" {
		if p, err := ioutil.ReadFile(src); err == nil {
			overview = template.HTML(string(p))
		}
	}
	return overview
}

func generateFuncMap(meta Meta) template.FuncMap {
	funcs := template.FuncMap{}
	funcs["baseURL"] = func() string {
		return meta.BaseURL
	}
	funcs["headers"] = func() []string {
		return meta.Headers
	}
	return funcs
}
