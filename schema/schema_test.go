package schema

import (
	"encoding/json"
	"reflect"
	"testing"
)

func readTestData(t *testing.T) (*Schema, *Schema) {
	user, err := newTestSchema(userJSON)
	if err != nil {
		t.Fatal(err)
	}
	article, err := newTestSchema(articleJSON)
	if err != nil {
		t.Fatal(err)
	}
	return user, article
}

func TestNewSchema(t *testing.T) {
	user, article := readTestData(t)

	s := article.ExampleJSON()
	var i map[string]interface{}
	if err := json.Unmarshal([]byte(s), &i); err != nil {
		t.Fatal("invalid article json")
	}
	articleUser, ok := i["user"].(map[string]interface{})
	if !ok {
		t.Fatal("invalid article json")

	}
	if !reflect.DeepEqual(articleUser, user.ExampleInterface()) {
		t.Errorf("article.user %v wants to equal user %v", articleUser, user.ExampleInterface())
	}
}

func TestExampleGetData(t *testing.T) {
	_, article := readTestData(t)

	for _, l := range article.Links {
		if l.Rel == "instances" && l.Href == "/articles" {
			getdata := l.Schema.ExampleGetData()
			if getdata[0] != "name=Article" {
				t.Errorf("expected name=Article, get %v", getdata[0])
			}
			return
		}
	}
	t.Errorf("instances link is not found")
}

func newTestSchema(str string) (*Schema, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, err
	}
	return NewSchema(data, nil)
}

const articleJSON = `{
  "$schema": "http://json-schema.org/draft-04/hyper-schema",
  "definitions": {
    "created_at": {
      "description": "when articles was created",
      "example": "2012-01-01T12:00:00Z",
      "format": "date-time",
      "type": [
        "string"
      ]
    },
    "id": {
      "description": "unique identifier of articles",
      "example": "01234567-89ab-cdef-0123-456789abcdef",
      "format": "uuid",
      "type": [
        "string"
      ]
    },
    "updated_at": {
      "description": "when articles was updated",
      "example": "2012-01-02T12:00:00Z",
      "format": "date-time",
      "type": [
        "string"
      ]
    }
  },
  "description": "FIXME",
  "id": "articles",
  "links": [
    {
      "description": "Create a new articles.",
      "href": "/articles",
      "method": "POST",
      "rel": "create",
      "schema": {
        "properties": {
          "id": {
            "$ref": "#/definitions/id"
          }
        },
        "type": [
          "object"
        ]
      },
      "title": "Create"
    },
    {
      "description": "Delete an existing articles.",
      "href": "/articles/{id}",
      "method": "DELETE",
      "rel": "destroy",
      "title": "Delete"
    },
    {
      "description": "Info for existing articles.",
      "href": "/articles/{id}",
      "method": "GET",
      "rel": "self",
      "title": "Info"
    },
    {
      "description": "List existing articles.",
      "href": "/articles",
      "method": "GET",
      "rel": "instances",
      "title": "List",
			"encType": "aplication/x-www-form-urlencoded",
			"schema": {
				"type": "object",
				"properties": {
					"name": {
						"descritption": "name of the product",
						"example": "Article"
					}
				}
			}
    },
    {
      "description": "Update an existing articles.",
      "href": "/articles/{id}",
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
    "created_at": {
      "$ref": "#/definitions/created_at"
    },
    "id": {
      "$ref": "#/definitions/id"
    },
    "updated_at": {
      "$ref": "#/definitions/updated_at"
    },
    "user": {
      "$ref": "user.json"
    }
  },
  "title": "FIXME - Articles",
  "type": [
    "object"
  ]
}
`

const userJSON = `{
  "$schema": "http://json-schema.org/draft-04/hyper-schema",
  "definitions": {
    "age": {
      "example": 16,
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
