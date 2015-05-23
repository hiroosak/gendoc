package schema

import "github.com/xeipuuv/gojsonschema"

// ValidSchema returns nil if json is valid.
func ValidSchema(jsonBytes []byte) error {
	schema := gojsonschema.NewStringLoader(string(jsonBytes))
	if _, err := gojsonschema.NewSchema(schema); err != nil {
		return err
	}
	return nil
}
