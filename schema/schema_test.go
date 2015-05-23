package schema

import (
	"encoding/json"
	"testing"
)

func TestNewSchema(t *testing.T) {
	user, err := newTestSchema(userJSON)
	if err != nil {
		t.Fatal("error")
	}
	article, err := newTestSchema(articleJSON)
	if err != nil {
		t.Fatal("error")
	}

	js := article.ExampleJSON()
	var i map[string]interface{}
	if err := json.Unmarshal([]byte(js), &i); err != nil {
		t.Fatal("error")
	}

	if i["created_at"] != "2012-01-01T12:00:00Z" {
		t.Error("wrong created_at")
	}
	if i["updated_at"] != "2012-01-01T12:00:00Z" {
		t.Error("wrong updated_at")
	}

	u1, ok := i["user"].(map[string]interface{})
	if !ok {
		t.Fatal("wrong user map")
	}

	js2 := user.ExampleJSON()
	var u2 map[string]interface{}
	if err := json.Unmarshal([]byte(js2), &u2); err != nil {
		t.Fatal("error")
	}
	keys := []string{"id", "age", "name"}
	for _, key := range keys {
		if u1[key] != u2[key] {
			t.Errorf("wrong %v", key)
		}
	}
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
      "example": "2012-01-01T12:00:00Z",
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
      "href": "/articless",
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
      "href": "/articless/{id}",
      "method": "DELETE",
      "rel": "destroy",
      "title": "Delete"
    },
    {
      "description": "Info for existing articles.",
      "href": "/articless/{id}",
      "method": "GET",
      "rel": "self",
      "title": "Info"
    },
    {
      "description": "List existing articless.",
      "href": "/articless",
      "method": "GET",
      "rel": "instances",
      "title": "List"
    },
    {
      "description": "Update an existing articles.",
      "href": "/articless/{id}",
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
