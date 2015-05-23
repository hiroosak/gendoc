package schema

import "testing"

func TestValidSchema(t *testing.T) {
	oks := []string{
		`{}`,
		`{"properties": {}}`,
	}

	for _, schema := range oks {
		if err := ValidSchema([]byte(schema)); err != nil {
			t.Errorf("%v validation is expected true. but error %v", schema, err)
		}
	}

	ngs := []string{
		`{"properties": []}`,
	}

	for _, schema := range ngs {
		if err := ValidSchema([]byte(schema)); err == nil {
			t.Errorf("%v validation is expected error. but not error", schema)
		}
	}
}
