package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var schemas map[string]*Schema = make(map[string]*Schema, 0)

type refPool struct {
	refMap map[string]*Schema
}

func NewRefPool() *refPool {
	return &refPool{
		refMap: make(map[string]*Schema, 0),
	}
}

func (r *refPool) Get(refStr string) *Schema {
	return r.refMap[refStr]
}

func (r *refPool) Set(refStr string, s *Schema) {
	r.refMap[refStr] = s
}

type Schema struct {
	Id          string
	Title       string
	Description string
	Type        string
	Format      string
	Example     interface{}
	Definitions map[string]*Schema
	Properties  map[string]*Schema

	Items []*Schema
	Links []LinkDescription

	Ref string

	CurrentRef string
	refPool    *refPool
	parent     *Schema
}

func NewSchemaFromInterface(data interface{}, refStr string, parent *Schema) (*Schema, error) {
	if refStr == "" {
		refStr = "#"
	}
	d, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("data type is not map[string]interface{}")
	}
	return NewSchema(d, refStr, parent)
}

func NewSchemaFromBytes(data []byte, refStr string, parent *Schema) (*Schema, error) {
	if refStr == "" {
		refStr = "#"
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, err
	}
	return NewSchema(dataMap, refStr, parent)
}

func NewSchema(data map[string]interface{}, refStr string, parent *Schema) (*Schema, error) {
	if refStr == "" {
		refStr = "#"
	}
	idStr := String(data, "id")
	typeStr := String(data, "type")
	s := &Schema{
		Type:        typeStr,
		Id:          String(data, "id"),
		Description: String(data, "description"),
		Format:      String(data, "format"),
		Title:       String(data, "title"),
		Example:     Interface(data, "example", typeStr),
		Ref:         String(data, "$ref"),
		CurrentRef:  refStr,
		parent:      parent,
	}
	s.Properties = make(map[string]*Schema, 0)
	s.Definitions = make(map[string]*Schema, 0)
	s.Items = make([]*Schema, 0)

	if idStr != "" {
		schemas[idStr] = s
	}

	// reference pool
	if parent != nil {
		s.refPool = parent.refPool
	}
	if s.refPool == nil {
		s.refPool = NewRefPool()
	}

	s.parseProperties(data["properties"])
	s.parseDefinitions(data["definitions"])
	s.parseItems(data["items"])

	s.refPool.Set(refStr, s)

	return s, nil
}

func (s *Schema) parseProperties(data interface{}) {
	properties, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	for key, property := range properties {
		prop, err := NewSchemaFromInterface(property, s.appendRefPath("properties", key), s)
		if err != nil {
			continue
		}
		s.Properties[key] = prop
	}
}

func (s *Schema) parseDefinitions(data interface{}) {
	definitions, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	for key, definition := range definitions {
		def, err := NewSchemaFromInterface(definition, s.appendRefPath("definitions", key), s)
		if err != nil {
			continue
		}
		s.Definitions[key] = def
	}
}

func (s *Schema) parseItems(data interface{}) error {
	item, err := NewSchemaFromInterface(data, s.appendRefPath("items"), s)
	if err != nil {
		return err
	}
	s.Items = append(s.Items, item)
	return nil
}

func (s *Schema) resolveReference(idStr, refStr string) *Schema {
	if refStr == "#" {
		return s
	}
	if schema, ok := schemas[idStr]; !ok {
		return s.refPool.Get(refStr)
	} else {
		return schema.refPool.Get(refStr)
	}
}

func (s *Schema) ExampleJSON() string {
	j := s.ExampleInterface()
	res, _ := json.MarshalIndent(j, "", "  ")
	return string(res)
}

func (s *Schema) ExampleInterface() interface{} {
	j := map[string]interface{}{}

	if example := fmt.Sprintf("%v", s.Example); example != "" {
		return s.Example
	}

	if s.Ref != "" {
		refs := s.resolveReference(s.Id, s.Ref)
		if refs != nil {
			return refs.ExampleInterface()
		}
	}

	if s.Type == "array" {
		return []interface{}{s.Items[0].ExampleInterface()}
	}

	for key, property := range s.Properties {

		if property.Ref != "" {
			refs := s.resolveReference(s.Id, property.Ref)
			j[key] = refs.ExampleInterface()
			continue
		}

		switch property.Type {
		//case "array":
		//	if len(property.Items) > 0 {
		//		j[key] = []interface{}{property.Items[0].ExampleInterface()}
		//	}
		case "object":
			j[key] = property.ExampleInterface()
		default:
			example := property.Example
			if example == "" {
				if _, ok := property.Properties[key]; ok {
					example = property.Properties[key].Example
				}
			}
			j[key] = example
		}
	}
	return j
}

func (s *Schema) appendRefPath(path ...string) string {
	paths := []string{s.CurrentRef}
	paths = append(paths, path...)

	return strings.Join(paths, "/")
}

type LinkDescription struct {
	Href         string
	Rel          string
	Title        string
	Description  string
	TargetSchema Schema
	MediaType    string
	Method       string
	EncType      string
	Schema       Schema
}
