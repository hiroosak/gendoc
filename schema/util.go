package schema

import "fmt"

func String(target interface{}, key string) string {
	d, ok := target.(map[string]interface{})
	if !ok {
		return ""
	}

	v, ok := d[key]
	if !ok {
		return ""
	}

	var res string

	switch v.(type) {
	case string:
		res = v.(string)
	case float64:
		res = fmt.Sprintf("%v", v)
	case []interface{}:
		rs := v.([]interface{})
		for _, r := range rs {
			if vv, ok := r.(string); ok {
				res = res + vv
			}
		}
	}
	return res
}

func Interface(target interface{}, key string, typeStr string) interface{} {
	if typeStr == "string" {
		return String(target, key)
	}

	var i interface{}
	d, ok := target.(map[string]interface{})
	if !ok {
		return i
	}

	v, ok := d[key]
	if !ok {
		return String(target, key)
	}
	return v
}
