package swagger

import (
	"encoding/json"
	"testing"
)

func TestJSONSchemaParse(t *testing.T) {
	var s JSONSchema
	if err := json.Unmarshal([]byte(`{"swagger": "2.0"}`), &s); err != nil {
		t.Fatal(err.Error())
	}
}
