package swagger

import (
	"github.com/hiroosak/gendoc/schema"
)

type JSONSchema struct {
	Swagger             string                `yaml:"swagger" json:"swagger"`
	Info                Info                  `yaml:"info,omitempty" json:"info,omitempty"`
	Host                string                `yaml:"host" json:"host"`
	BasePath            String                `yaml:"basePath,omitempty" json:"basePath,omitempty"`
	Schemes             Schemes               `yaml:"schemes,omitempty" json:"schemes,omitempty"`
	Consumes            *Consumes             `yaml:"consumes,omitempty" json:"consumes,omitempty"`
	Produces            Produces              `yaml:"produces,omitempty" json:"produces,omitempty"`
	Paths               Paths                 `yaml:"paths,omitempty" json:"paths,omitempty"`
	Definitions         Definitions           `yaml:"definitions,omitempty" json:"definitions,omitempty"`
	Parameters          *ParameterDefinitions `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses           ResponseDefinitions   `yaml:"responses,omitempty" json:"responses,omitempty"`
	Security            *Security             `yaml:"security,omitempty" json:"security,omitempty"`
	SecurityDefinitions *SecurityDefinitions  `yaml:"securityDefinitions,omitempty" json:"securityDefinitions,omitempty"`
	Tags                *Tags                 `yaml:"tags,omitempty" json:"tags,omitempty"`
	ExternalDocs        *ExternalDocs         `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

type String *string

func NewString(s string) *string {
	return &s
}

type Schemes []Scheme
type Definitions map[string]DefinitionsParams
type DefinitionsParams map[string]interface{}
type Paths map[string]*Path
type Scheme string

type Consumes MediaTypeList

type Produces MediaTypeList

func NewProduces(ss []string) Produces {
	ps := make(Produces, len(ss))
	for k, s := range ss {
		ps[k] = MediaType(s)
	}
	return ps
}

type MediaTypeList []MediaType
type ParameterDefinitions []ParameterDefinition
type ResponseDefinitions map[string]ResponseDefinition
type SecurityDefinitions []SecurityDefinition
type Tags []Tag
type ExternalDocs []ExternalDoc
type Responses []Response

type MediaType string

type ResponseDefinition struct {
	Description string        `yaml:"description,omitempty" json:"description,omitempty"`
	Schema      schema.Schema `yaml:"schema,omitempty" json:"schema,omitempty"`
	Headers     Headers       `yaml:"headers,omitempty" json:"headers,omitempty"`
	Ref         string        `yaml:"$ref,omitempty" json:"$ref,omitempty"`
}

type ExternalDoc struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	URL         string `yaml:"url" json:"url"`
}

type SecurityDefinition interface{}

type Contact struct {
	Name  string  `yaml:"name" json:"name"`
	Url   *string `yaml:"url,omitempty" json:"url,omitempty"`
	Email *string `yaml:"email,omitempty" json:"email,omitempty"`
}

type License struct {
	Name string `yaml:"string" json:"string"`
	Url  string `yaml:"url" json:"url"`
}

type Path struct {
	Ref        string     `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Get        *Operation `yaml:"get,omitempty" json:"get,omitempty"`
	Put        *Operation `yaml:"put,omitempty" json:"put,omitempty"`
	Post       *Operation `yaml:"post,omitempty" json:"post,omitempty"`
	Delete     *Operation `yaml:"delete,omitempty" json:"delete,omitempty"`
	Options    *Operation `yaml:"options,omitempty" json:"options,omitempty"`
	Head       *Operation `yaml:"head,omitempty" json:"head,omitempty"`
	Patch      *Operation `yaml:"patch,omitempty" json:"patch,omitempty"`
	Parameters Parameters `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

type Parameters []ParameterDefinitions

type ParameterDefinition struct {
	In               string             `yaml:"in" json:"in"`
	Name             string             `yaml:"name" json:"name"`
	Description      string             `yaml:"description,omitempty" json:"description,omitempty"`
	Required         bool               `yaml:"required" json:"required"`
	Type             string             `yaml:"type,omitempty" json:"type,omitempty"`
	Schema           *schema.Schema     `yaml:"schema,omitempty" json:"schema,omitempty"`
	CollectionFormat []CollectionFormat `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Default          string             `yaml:"default,omitempty" json:"default,omitempty"`
	Maximum          int                `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	ExclusiveMaximum bool               `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	Minimum          int                `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	ExclusiveMinimum bool               `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	MaxLength        int                `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	MinLength        int                `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxItems         int                `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	MinItems         int                `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	UniqueItems      bool               `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`
	Enum             string             `yaml:"enum,omitempty" json:"enum,omitempty"`
}

type Operation struct {
	Tags         []string            `yaml:"tags,omitempty" json:"tags,omitempty"`
	Summary      string              `yaml:"summary" json:"summary"`
	Description  string              `yaml:"description" json:"description"`
	ExternalDocs *ExternalDocs       `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	OperationId  string              `yaml:"operationId" json:"operationId"`
	Produces     Produces            `yaml:"produces,omitempty" json:"produces,omitempty"`
	Consumes     Consumes            `yaml:"consumes,omitempty" json:"consumes,omitempty"`
	Parameters   Parameters          `yaml:"parameters" json:"parameters"`
	Response     map[string]Response `yaml:"responses" json:"responses"`
	Schemes      *Schemes            `yaml:"schemes,omitempty" json:"schemes,omitempty"`
}

type Response struct {
	Description string                 `yaml:"description" json:"description" `
	Schema      map[string]interface{} `yaml:"schema,omitempty" json:"schema,omitempty" `
	Headers     Headers                `yaml:"headers,omitempty" json:"headers,omitempty" `
	Examples    Examples               `yaml:"examples,omitempty" json:"examples,omitempty" `
}

type Examples interface{}

type Headers map[string]Header

type Header struct {
	Type   string          `yaml:"type" json:"type"`
	Format string          `yaml:"string,omitempty" json:"string,omitempty"`
	Item   PrimitivesItems `yaml:"items,omitempty" json:"items,omitempty"`
}

type PrimitivesItems struct {
	Type             string             `yaml:"type" json:"type"`
	Format           string             `yaml:"string,omitempty" json:"string,omitempty"`
	Items            interface{}        `yaml:"items,omitempty" json:"items,omitempty"`
	CollectionFormat []CollectionFormat `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Default          string             `yaml:"default,omitempty" json:"default,omitempty" `
	Maximum          int                `yaml:"maximum,omitempty" json:"maximum,omitempty" `
	ExclusiveMaximum bool               `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	Minimum          int                `yaml:"minimum,omitempty" json:"minimum,omitempty" `
	ExclusiveMinimum bool               `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	MaxLength        int                `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	MinLength        int                `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxItems         int                `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	MinItems         int                `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	UniqueItems      bool               `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`
}

type CollectionFormat string

type Security []SecurityRequirement
type SecurityRequirement []string

type XML struct {
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
	Prefix    string `yaml:"prefix" json:"prefix"`
	Attribute bool   `yaml:"attribute" json:"attribute"`
	Wrapped   bool   `yaml:"wrapped" json:"wrapped"`
}

type Tag struct {
	Name         string        `yaml:"name" json:"name"`
	Description  String        `yaml:"description,omitempty" json:"description,omitempty"`
	externalDocs *ExternalDocs `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

type Info struct {
	Title          string   `yaml:"title" json:"title"`
	Version        String   `yaml:"version" json:"version"`
	Description    String   `yaml:"description,omitempty" json:"description,omitempty"`
	TermsOfService String   `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	Contact        *Contact `yaml:"contact,omitempty" json:"contact,omitempty"`
	License        *License `yaml:"license,omitempty" json:"license,omitempty"`
}
