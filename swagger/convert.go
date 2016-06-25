package swagger

import (
	"fmt"
	"io"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"sync"

	_ "github.com/ghodss/yaml"
	"gopkg.in/yaml.v2"

	"github.com/hiroosak/gendoc/schema"
)

type Param struct {
	Title       string
	Description string
	BaseURL     string
	Version     string
	Produces    []string
	Schemes     Schemes
	Headers     []string
}

func (p Param) Host() string {
	u, _ := url.Parse(p.BaseURL)
	return u.Host
}

func (p Param) Path() string {
	u, _ := url.Parse(p.BaseURL)
	return u.Path
}

type Converter struct {
	mutex sync.Mutex

	schema *JSONSchema
}

func NewConverter(p Param) *Converter {
	return &Converter{
		mutex: sync.Mutex{},
		schema: &JSONSchema{
			Swagger: "2.0",
			Info: Info{
				Title:       p.Title,
				Version:     NewString(p.Version),
				Description: NewString(p.Description),
			},
			Schemes:     p.Schemes,
			Host:        p.Host(),
			BasePath:    NewString((p.Path())),
			Produces:    NewProduces(p.Produces),
			Paths:       map[string]*Path{},
			Definitions: Definitions{},
		},
	}
}

func (c *Converter) Convert(w io.Writer, ss schema.SchemaSlice) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, s := range ss {
		title := strings.ToTitle(string(s.Id[0]))
		resource := title + s.Id[1:]

		properties := Definitions{}

		for key, property := range s.Properties {
			typ := property.ResolveType()
			types := ""
			if len(typ) > 0 {
				types = typ[0]
			}
			if ss := strings.Split(types, " "); len(ss) > 1 {
				types = ss[0]
			}

			examples := property.ExampleJSON()

			ty := "string"
			switch reflect.ValueOf(property.ExampleInterface()).Kind() {
			case reflect.String:
				ty = "string"
			case reflect.Float64:
				ty = "integer"
			case reflect.Slice:
				ty = "array"
			case reflect.Map:
				ty = "string"
			}

			properties[string(key)] = DefinitionsParams{
				"description": property.ResolveDescription(),
				"example":     examples,
				"type":        ty,
			}
			if ty == "array" {
				properties[string(key)]["items"] = map[string]string{
					"type": "string",
				}
			}
		}

		definitions := DefinitionsParams{
			"type":       "object",
			"properties": properties,
		}

		c.schema.Definitions[string(resource)] = definitions

		for _, link := range s.Links {
			if _, ok := c.schema.Paths[link.Href]; !ok {
				c.schema.Paths[link.Href] = &Path{}
			}
			parameters := Parameters{}
			for key, property := range link.Schema.Properties {
				typs := property.ResolveType()
				typ := "string"
				if len(typs) > 1 {
					typ = typs[0]
				}

				var in string
				switch link.Method {
				case "GET", "Get", "get":
					in = "query"
				case "POST", "Post", "post":
					in = "body"
				case "PUT", "Put", "put":
					in = "body"
				case "DELETE", "Delete", "delete":
					in = "query"
				case "PATCH", "Patch", "patch":
					in = "body"
				}

				param := ParameterDefinitions{
					ParameterDefinition{
						Name:        key,
						Description: property.Description,
						Type:        typ,
						In:          in,
					},
				}
				matched, err := regexp.MatchString(fmt.Sprintf("{%v}", key), link.Href)
				if err == nil && matched {
					param[0].In = "path"
					param[0].Required = true
				}
				if in == "query" {
					param[0].Type = "string"
				}

				parameters = append(parameters, param)
			}

			ope := &Operation{
				Summary:     link.Title,
				Description: link.Description,
				Parameters:  parameters,
				OperationId: fmt.Sprintf("%v-%v-%v", link.Method, string(resource), link.Title),
				Consumes:    Consumes(MediaTypeList([]MediaType{"application/json"})),
				Produces:    Produces(MediaTypeList([]MediaType{"application/json"})),
				Response:    map[string]Response{},
			}
			switch link.Method {
			case "GET", "Get", "get":
				ope.Response = map[string]Response{
					"200": Response{
						Description: link.Description,
						Schema: map[string]interface{}{
							"example": link.TargetSchema.ExampleJSON(),
						},
					},
				}
				h := c.schema.Paths[link.Href]
				h.Get = ope
			case "POST", "Post", "post":
				ope.Response = map[string]Response{
					"201": Response{
						Description: link.Description,
					},
				}
				ps := []ParameterDefinition{
					ParameterDefinition{
						Name:        "body",
						Description: link.Description,
						In:          "body",
						Schema: &schema.Schema{
							Ref: fmt.Sprintf("#/definitions/%s", string(resource)),
						},
					},
				}
				ope.Parameters = []ParameterDefinitions{ps}
				h := c.schema.Paths[link.Href]
				h.Post = ope
			case "PUT", "Put", "put":
				ope.Response = map[string]Response{
					"200": Response{
						Description: link.Description,
					},
				}
				ps := []ParameterDefinition{
					ParameterDefinition{
						Name:        "body",
						Description: link.Description,
						In:          "body",
						Schema: &schema.Schema{
							Ref: fmt.Sprintf("#/definitions/%s", string(resource)),
						},
					},
				}
				ope.Parameters = []ParameterDefinitions{ps}
				h := c.schema.Paths[link.Href]
				h.Put = ope
			case "DELETE", "Delete", "delete":
				ope.Response = map[string]Response{
					"200": Response{
						Description: link.Description,
					},
				}
				h := c.schema.Paths[link.Href]
				h.Delete = ope
			case "PATCH", "Patch", "patch":
				ope.Response = map[string]Response{
					"203": Response{
						Description: link.Description,
					},
				}
				ps := []ParameterDefinition{
					ParameterDefinition{
						Name:        "body",
						Description: link.Description,
						In:          "body",
						Schema: &schema.Schema{
							Ref: fmt.Sprintf("#/definitions/%s", string(resource)),
						},
					},
				}
				ope.Parameters = []ParameterDefinitions{ps}
				h := c.schema.Paths[link.Href]
				h.Patch = ope
			}
		}
	}

	bs, err := yaml.Marshal(c.schema)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}
