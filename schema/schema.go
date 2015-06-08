package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
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

func (r *refPool) Keys() []string {
	rs := make([]string, r.Len())
	var i int
	for k, _ := range r.refMap {
		rs[i] = k
		i += 1
	}
	return rs
}

func (r *refPool) Len() int {
	return len(r.refMap)
}

// +gen slice:"GroupBy[string],Where"
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
	Links []*LinkDescription

	Ref string

	CurrentRef string
	refPool    *refPool
	parent     *Schema
}

func NewSchemaFromFile(path string, info os.FileInfo) (*Schema, error) {
	bytes, err := YamlFileToJson(path, info)
	if err != nil {
		return nil, err
	}
	return NewSchemaFromBytes(bytes, "", nil)
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
	s.Links = make([]*LinkDescription, 0)

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
	s.parseLinks(data["links"])
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

func (s *Schema) parseLinks(data interface{}) error {
	if data == nil {
		return nil
	}

	linkLists, ok := data.([]interface{})
	if !ok {
		return errors.New("parse failed links")
	}

	for i, l := range linkLists {
		link, ok := l.(map[string]interface{})
		if !ok {
			return errors.New("parse failed links")
		}
		var schema *Schema
		if v, ok := link["schema"]; ok {
			schema, _ = NewSchemaFromInterface(v, s.appendRefPath(fmt.Sprintf("links[%v]", i), "schema"), s)
		} else {
			schema = s
		}
		var targetSchema *Schema
		if v, ok := link["targetSchema"]; ok {
			targetSchema, _ = NewSchemaFromInterface(v, s.appendRefPath(fmt.Sprintf("links[%v]", i), "targetSchema"), s)
		} else {
			targetSchema = s
		}

		l := &LinkDescription{
			Title:        String(link, "title"),
			Description:  String(link, "description"),
			Href:         String(link, "href"),
			Method:       String(link, "method"),
			Rel:          String(link, "rel"),
			EncType:      String(link, "encType"),
			Schema:       schema,
			TargetSchema: targetSchema,
		}

		s.Links = append(s.Links, l)
	}

	return nil
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
	idStr, refStr = parseReference(idStr, refStr)
	if schema, ok := schemas[idStr]; !ok {
		return s.refPool.Get(refStr)
	} else {
		return schema.refPool.Get(refStr)
	}
}

func parseReference(idStr, refStr string) (string, string) {
	if n := strings.Index(refStr, "."); n > 0 {
		idStr = refStr[0:n]
		if h := strings.Index(refStr, "#"); h > 0 {
			refStr = refStr[h:len(refStr)]
		} else {
			refStr = "#"
		}
	}
	return idStr, refStr
}

func (s *Schema) Alias() *Schema {
	if s.Ref == "" {
		return s
	}
	return s.resolveReference(s.Id, s.Ref)
}

func (s *Schema) ResolveType() string {
	schema := s.Alias()
	if schema == nil {
		return ""
	}
	return schema.Type
}

func (s *Schema) ResolveFormat() string {
	schema := s.Alias()
	if schema == nil {
		return ""
	}
	return schema.Format
}

func (s *Schema) ResolveDescription() string {
	schema := s.Alias()
	if schema == nil {
		return ""
	}
	return schema.Description
}

func (s *Schema) ExampleJSON() string {
	j := s.ExampleInterface()
	res, _ := json.MarshalIndent(j, "", "  ")
	return string(res)
}

func (s *Schema) ExampleInterface() interface{} {
	if s == nil {
		return nil
	}

	if s.Example != nil {
		if example := fmt.Sprintf("%v", s.Example); example != "" {
			return s.Example
		}
	}

	if s.Ref != "" {
		if refs := s.resolveReference(s.Id, s.Ref); refs != nil {
			return refs.ExampleInterface()
		}
	}

	if s.Type == "array" {
		return []interface{}{s.Items[0].ExampleInterface()}
	}

	j := map[string]interface{}{}
	for key, property := range s.Properties {
		if property.Ref != "" {
			refs := s.resolveReference(s.Id, property.Ref)
			j[key] = refs.ExampleInterface()
		} else {
			j[key] = property.Example
		}
	}
	return j
}

func (s *Schema) ExampleGetData() []string {
	p := s.ExampleInterface()
	params, ok := p.(map[string]interface{})
	if !ok {
		return []string{}
	}
	val := url.Values{}
	for k, v := range params {
		val.Set(k, fmt.Sprintf("%v", v))
	}
	e := val.Encode()
	return strings.Split(e, "&")
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
	MediaType    string
	Method       string
	EncType      string
	Schema       *Schema
	TargetSchema *Schema
}
