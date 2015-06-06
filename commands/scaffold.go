package commands

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"bitbucket.org/pkg/inflect"
)

func Scaffold(r string) {
	if r == "" {
		fmt.Println("Please input resource name")
		fmt.Println("USAGE:")
		fmt.Println("   gendoc init article > src/article.yml")
		fmt.Println("")
		return
	}
	w := renderScaffold(r)
	w.WriteTo(os.Stdout)
}

func renderScaffold(r string) *bytes.Buffer {
	resource := inflect.Singularize(r)
	resources := inflect.Pluralize(resource)

	tmpl := template.Must(template.New("scaffold").Parse(schemaTmpl))

	w := bytes.NewBuffer([]byte{})
	p := struct {
		Resource  string
		Resources string
	}{
		Resource:  resource,
		Resources: resources,
	}
	if err := tmpl.Execute(w, p); err != nil {
		panic(err)
	}
	return w
}

const schemaTmpl = `---
$schema: "http://json-schema.org/draft-04/hyper-schema"
id: "{{ .Resource }}"
title: {{ .Resource }}
type: object
definitions:
  id:
    type: integer
    description: resource id
    example: 1
  name:
    type: string
    description: user name
    example: Ken
  createdAt:
    type: string
    description: datetime created data
    format: "date-time"
    example: 2015-04-21T23:59:60Z
  updatedAt:
    type: string
    description: datetime updated data
    format: "date-time"
    example: 2015-04-21T23:59:60Z
links:
- title: List
  description: List existing {{ .Resources }}.
  href: "/{{ .Resources }}"
  method: GET
  rel: instances
- title: Info
  description: Info for existing {{ .Resource }}.
  href: "/{{ .Resources }}/{id}"
  method: GET
  rel: self
- title: Create
  description: Create a new {{ .Resource }}.
  href: "/{{ .Resources }}"
  method: POST
  rel: create
  schema:
    properties: {}
    type:
    - object
- title: Update
  description: Update an existing {{ .Resource }}.
  href: "/{{ .Resources }}/{id}"
  method: PATCH
  rel: update
  schema:
    properties: {}
    type:
    - object
- title: Delete
  description: Delete an existing {{ .Resource }}.
  href: "/{{ .Resource }}/{id}"
  method: DELETE
  rel: destroy
properties:
  id:
    "$ref": "#/definitions/id"
  name:
    "$ref": "#/definitions/name"
  createdAt:
    "$ref": "#/definitions/createdAt"
  updatedAt:
    "$ref": "#/definitions/updatedAt"
required:
- id 
- name
- createdAt
- updatedAt`
