package schema

import (
	"encoding/json"
	"testing"
)

func TestExampleJSON(t *testing.T) {
	s, err := NewSchemaFromBytes([]byte(userJSON), "", nil)
	if err != nil {
		t.Fatal("invalid json")
	}
	jsonStr := s.ExampleJSON()

	var i map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &i); err != nil {
		t.Fatalf("invalid example json. %v", err)
	}
	if i["age"].(float64) != 16 {
		t.Errorf("age is expected 16. but %v", i["age"])
	}
	if i["id"].(string) != "1" {
		t.Errorf("id is expected 1. but %v", i["id"])
	}
	if i["name"].(string) != "Mike" {
		t.Errorf("name is expected Mike. but %v", i["name"])
	}
}

var t1 *Schema = &Schema{
	Id: "user",
	Definitions: map[string]*Schema{
		"username": &Schema{
			Description: "user's name",
			Type:        "string",
			Example:     "Ken",
		},
		"age": &Schema{
			Description: "user's old",
			Type:        "integer",
			Example:     18,
		},
	},
	Type: "array",
	Items: []*Schema{
		&Schema{
			Ref: "#",
		},
	},
}

func TestExampleJSONArray(t *testing.T) {
	var jsonstr = `{
		"id": "user",
		"definitions": {
			"age": {
				"example": 16,
				"description": "user's age",
				"type": "integer"
			},
			"id": {
				"example": 1,
				"type": "string"
			}
		},
		"type": "array",
		"items": {
			"$ref": "#/definitions/age"
		}
	}`

	t1, err := NewSchemaFromBytes([]byte(jsonstr), "", nil)
	if err != nil {
		t.Fatal(err)
	}

	j := t1.ExampleJSON()
	var jm []interface{}
	json.Unmarshal([]byte(j), &jm)

	if len(jm) != 1 {
		t.Fatalf("len(j) length is not 1. %v", len(jm))
	}
	if jm[0].(float64) != 16 {
		t.Errorf("age is not 16. %v", j)
	}
}

func TestExampleJSONArrayObject(t *testing.T) {
	var jsonstr = `{
		"id": "user",
		"definitions": {
			"data": {
				"type": "object",
				"properties": {
					"age": {
						"example": 16,
						"description": "user's age",
						"type": "integer"
					}
				}
			}
		},
		"type": "array",
		"items": {
			"$ref": "#/definitions/data"
		}
	}`

	t1, err := NewSchemaFromBytes([]byte(jsonstr), "", nil)
	if err != nil {
		t.Fatal(err)
	}

	j := t1.ExampleJSON()
	var jm []map[string]interface{}
	err = json.Unmarshal([]byte(j), &jm)
	if err != nil {
		t.Fatal(err)
	}
	if len(jm) != 1 {
		t.Fatalf("len(j) length is not 1. %v", len(jm))
	}
	if jm[0]["age"].(float64) != 16 {
		t.Errorf("age is not 16. %v", j)
	}
}

func TestRefpool(t *testing.T) {
	_, err := NewSchemaFromBytes([]byte(userJSON), "", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestResolveReference(t *testing.T) {
	s, err := NewSchemaFromBytes([]byte(userJSON), "", nil)
	if err != nil {
		t.Fatal("invalid json")
	}
	r1 := s.resolveReference("user", "#/definitions/age")
	if r1.Description != "user's age" {
		t.Errorf("description is not equal")
	}
	if r1.Example != float64(16) {
		t.Errorf("example is not equal %v != %v", r1.Example, 16)
	}
}

func TestAppendRefPath(t *testing.T) {
	type t2 struct {
		src []string
		dst string
	}

	res := []t2{
		t2{[]string{"a"}, "#/a"},
		t2{[]string{"definitions", "id"}, "#/definitions/id"},
	}

	s := &Schema{CurrentRef: "#"}
	for _, r := range res {
		if s.appendRefPath(r.src...) != r.dst {
			t.Errorf("expected %v. but actual %v", r.src, r.dst)
		}
	}
}

func TestExampleGetData(t *testing.T) {
	accepts := []string{
		"age=16",
		"id=1",
		"name=Mike",
	}

	s, _ := NewSchemaFromBytes([]byte(userJSON), "", nil)
	datas := s.ExampleGetData()

	for k, v := range datas {
		if v != accepts[k] {
			t.Errorf("accept %v, but TV", v, accepts[k])
		}
	}
}

func TestParseReference(t *testing.T) {
	type t1 struct {
		SrcId  string
		SrcRef string
		DstId  string
		DstRef string
	}
	res := []t1{
		t1{"", "user.json", "user", "#"},
		t1{"", "user.json#definitions/id", "user", "#definitions/id"},
		t1{"user", "", "user", ""},
	}
	for _, r := range res {
		id, ref := parseReference(r.SrcId, r.SrcRef)
		if id != r.DstId {
			t.Errorf("%v: id is expected %v. but %v", r, r.DstId, id)
		}
		if ref != r.DstRef {
			t.Errorf("%v: ref is expected [%v]. but [%v]", r, r.DstRef, ref)
		}
	}
}

const userJSON = `{
  "$schema": "http://json-schema.org/draft-04/hyper-schema",
  "definitions": {
    "age": {
      "example": 16,
			"description": "user's age",
			"type": "integer"
    },
    "id": {
      "example": 1,
      "type": "string"
    },
    "name": {
      "example": "Mike",
      "type": "string"
    }
  },
  "id": "user",
  "links": [
    {
      "description": "Create a new user.",
      "href": "/users",
      "method": "POST",
      "rel": "create",
      "schema": {
        "properties": {},
        "type": [
          "object"
        ]
      },
      "title": "Create"
    },
    {
      "description": "Delete an existing user.",
      "href": "/users/{id}",
      "method": "DELETE",
      "rel": "destroy",
      "title": "Delete"
    },
    {
      "description": "Info for existing user.",
      "href": "/users/{id}",
      "method": "GET",
      "rel": "self",
      "title": "Info"
    },
    {
      "description": "List existing users.",
      "href": "/users",
      "method": "GET",
      "rel": "instances",
      "title": "List"
    },
    {
      "description": "Update an existing user.",
      "href": "/users/{id}",
      "method": "PATCH",
      "rel": "update",
      "schema": {
        "properties": {},
        "type": [
          "object"
        ]
      },
      "title": "Update"
    }
  ],
  "properties": {
    "age": {
      "$ref": "#/definitions/age"
    },
    "id": {
      "$ref": "#/definitions/id"
    },
    "name": {
      "$ref": "#/definitions/name"
    }
  },
  "required": [
    "id",
    "name",
    "age"
  ],
  "title": "user",
  "type": "object"
}`
