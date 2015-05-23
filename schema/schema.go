package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/mattn/go-scan"
)

const dummyType = "resource"

var allData map[string]Schema = map[string]Schema{}

// +gen slice:"GroupBy[string],Where"
type Schema struct {
	Id          string
	Title       string
	Description string
	Type        string
	Format      string
	Example     interface{}
	Definitions map[string]Schema
	Properties  map[string]Schema

	Items []Schema
	Links []LinkDescription

	data     map[string]interface{}
	rootData map[string]interface{}
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

func NewSchema(data, rootData map[string]interface{}) (*Schema, error) {
	if rootData == nil {
		rootData = data
	}
	typeStr := String(data, "type")
	r := &Schema{
		data:        data,
		rootData:    rootData,
		Type:        typeStr,
		Id:          String(data, "id"),
		Description: String(data, "descption"),
		Format:      String(data, "format"),
		Title:       String(data, "title"),
		Example:     Interface(data, "example", typeStr),
	}

	r.setItems()
	r.setDefinitions()
	r.setProperties()
	r.setLinkList()

	allData[r.Id] = *r

	return r, nil
}

func NewSchemaFromFile(path string, info os.FileInfo) (*Schema, error) {
	isJSON := isExtJSONFile(info)
	isYAML := isExtYaml(info)

	if !isJSON && !isYAML {
		return nil, fmt.Errorf("%v is not support file format", info.Name())
	}

	rs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var d map[string]interface{}
	switch {
	case isJSON:
		if err := json.Unmarshal(rs, &d); err != nil {
			return nil, err
		}
	case isYAML:
		if err := yaml.Unmarshal(rs, &d); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%v is not support file format", info.Name())
	}

	return NewSchema(d, nil)
}

func (s Schema) ExampleAlias() interface{} {
	if s.Type == dummyType {
		e := s.Example.(string)
		if a, ok := allData[baseResourceName(e)]; ok {
			return a.ExampleJSON()
		}
	}
	return s.Example
}

func NewSchemaFromInterface(data, rootData interface{}) (*Schema, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}
	d, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data type is not map[string]interface{}")
	}
	rd, ok := rootData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("rootData type is not map[string]interface{}")
	}

	return NewSchema(d, rd)
}

func (r *Schema) setDefinitions() error {
	if err := r.setRefDefinitions(); err == nil {
		return nil
	}
	var definitions map[string]interface{}
	if err := scan.ScanTree(r.data, `/definitions`, &definitions); err != nil {
		return err
	}

	r.Definitions = make(map[string]Schema, len(definitions))
	for name, d := range definitions {
		r.Definitions[name] = Schema{
			Description: String(d, "description"),
			Type:        String(d, "type"),
			Example:     Interface(d, "example", String(d, "type")),
			Format:      String(d, "format"),
		}
	}
	return nil
}

func (r *Schema) setRefDefinitions() error {
	var ref string
	if err := scan.ScanTree(r.data, "/$ref", &ref); err == nil {
		if isJSONExt(ref) {
			r.Definitions = make(map[string]Schema, 1)
			r.Definitions[baseResourceName(ref)] = Schema{
				Example: ref,
				Type:    dummyType,
			}
			return nil
		}

		i := strings.Index(ref, "#")
		path := ref[i+1 : len(ref)]

		var t map[string]interface{}
		if scan.ScanTree(r.rootData, path, &t); err == nil {
			ss := strings.Split(ref, "/")
			name := ss[len(ss)-1]
			r.Definitions = make(map[string]Schema, 1)
			r.Definitions[name] = Schema{
				Description: String(t, "description"),
				Type:        String(t, "type"),
				Example:     Interface(t, "example", String(t, "type")),
				Format:      String(t, "format"),
			}
		}
	}
	return nil
}

func (r *Schema) setLinkList() error {
	var linksData []map[string]interface{}
	if err := scan.ScanTree(r.data, `/links`, &linksData); err != nil {
		return err
	}

	for _, link := range linksData {
		schema, err := NewSchemaFromInterface(link["schema"], r.rootData)
		if err != nil || len(schema.Properties) == 0 {
			schema = r
		}
		targetSchema, err := NewSchemaFromInterface(link["targetSchema"], r.rootData)
		if err != nil || len(targetSchema.Properties) == 0 {
			targetSchema = r
		}

		l := LinkDescription{
			Title:        String(link, "title"),
			Description:  String(link, "description"),
			Href:         String(link, "href"),
			Method:       String(link, "method"),
			Rel:          String(link, "rel"),
			Schema:       *schema,
			TargetSchema: *targetSchema,
		}
		r.Links = append(r.Links, l)
	}

	return nil
}

func (r *Schema) setItems() error {
	var i interface{}
	if err := scan.ScanTree(r.data, `/items`, &i); err != nil {
		return err
	}

	items, err := NewSchemaFromInterface(i, r.rootData)
	if err != nil {
		return err
	}
	r.Items = append(r.Items, *items)

	return nil
}

func (r *Schema) setProperties() error {
	var propertyData map[string]interface{}
	if err := scan.ScanTree(r.data, `/properties`, &propertyData); err != nil {
		return err
	}

	r.Properties = make(map[string]Schema, len(propertyData))

	for name, property := range propertyData {
		var ref string
		if err := scan.ScanTree(property, "/$ref", &ref); err == nil {
			if isJSONExt(ref) {
				r.Properties[baseResourceName(name)] = Schema{
					Example: ref,
					Type:    dummyType,
				}
				continue
			}

			i := strings.Index(ref, "#")
			path := ref[i+1 : len(ref)]

			var t map[string]interface{}
			if scan.ScanTree(r.rootData, path, &t); err == nil {
				r.Properties[baseResourceName(ref)] = Schema{
					Description: String(t, "description"),
					Type:        String(t, "type"),
					Example:     Interface(t, "example", String(t, "type")),
					Format:      String(t, "format"),
				}
				continue
			}
		}
		if schema, err := NewSchemaFromInterface(property, r.rootData); err == nil {
			r.Properties[name] = *schema
		}
	}

	return nil
}

func (s *Schema) ExampleJSON() string {
	j := s.ExampleInterface()
	res, _ := json.MarshalIndent(j, "", "  ")
	return string(res)
}

func (s *Schema) ExampleInterface() map[string]interface{} {
	j := map[string]interface{}{}

	for key, s := range s.Properties {
		switch s.Type {
		case "array":
			if len(s.Items) > 0 {
				j[key] = []interface{}{s.Items[0].ExampleInterface()}
			}
		case "object":
			j[key] = s.ExampleInterface()
		default:
			if s.Type == dummyType {
				e := s.Example.(string)
				if a, ok := allData[baseResourceName(e)]; ok {
					j[key] = a.ExampleInterface()
				}
			} else {
				example := s.Example
				if example == "" {
					example = s.Properties[key].Example
				}
				j[key] = example
			}
		}
	}
	return j
}

func (s *Schema) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s.data, "", " ")
}
