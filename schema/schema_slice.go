// Generated by: gen
// TypeWriter: slice
// Directive: +gen on Schema

package schema

// SchemaSlice is a slice of type Schema. Use it where you would use []Schema.
type SchemaSlice []Schema

// GroupByString groups elements into a map keyed by string. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv SchemaSlice) GroupByString(fn func(Schema) string) map[string]SchemaSlice {
	result := make(map[string]SchemaSlice)
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Where returns a new SchemaSlice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv SchemaSlice) Where(fn func(Schema) bool) (result SchemaSlice) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}