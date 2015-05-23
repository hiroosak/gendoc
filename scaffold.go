package main

import (
	"fmt"
	"strings"

	"bitbucket.org/pkg/inflect"
)

func scaffold(r string) {
	if r == "" {
		fmt.Println("Please input resource name")
		fmt.Println("USAGE:")
		fmt.Println("   gendoc init article > src/article.yml")
		fmt.Println("")
		return
	}
	resource := inflect.Singularize(r)
	resources := inflect.Pluralize(resource)

	s := strings.Replace(schemaTmpl, "[resource_name]", resource, -1)
	s = strings.Replace(s, "[resource_names]", resources, -1)

	fmt.Println(s)
}

const schemaTmpl = `---
$schema: "http://json-schema.org/draft-04/hyper-schema"
id: "[resource_name]"
title: [resource_name]
type: object
definitions:
  id:
    type: integer
    example: 1
  name:
    type: string
    example: Ken
  createdAt:
    type: string
    format: "date-time"
    example: 2015-04-21T23:59:60Z
  updatedAt:
    type: string
    format: "date-time"
    example: 2015-04-21T23:59:60Z
links:
- description: List existing [resource_names].
  href: "/[resource_names]"
  method: GET
  rel: instances
  title: List
- description: Info for existing [resource_name].
  href: "/[resource_names]/{id}"
  method: GET
  rel: self
  title: Info
- description: Create a new [resource_name].
  href: "/[resource_names]"
  method: POST
  rel: create
  schema:
    properties: {}
    type:
    - object
  title: Create
- description: Update an existing [resource_name].
  href: "/[resource_names]/{id}"
  method: PATCH
  rel: update
  schema:
    properties: {}
    type:
    - object
  title: Update
- description: Delete an existing [resource_name].
  href: "/[resource_name]/{id}"
  method: DELETE
  rel: destroy
  title: Delete
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
