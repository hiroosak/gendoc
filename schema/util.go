package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

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

func StringSlice(target interface{}, key string) []string {
	d, ok := target.(map[string]interface{})
	if !ok {
		return []string{}
	}

	v, ok := d[key]
	if !ok {
		return []string{}
	}

	switch v.(type) {
	case string:
		return []string{v.(string)}
	case float64:
		return []string{fmt.Sprintf("%v", v)}
	case []interface{}:
		res := []string{}
		rs := v.([]interface{})
		for _, r := range rs {
			if vv, ok := r.(string); ok {
				res = append(res, vv)
			}
		}
		return res
	}
	return []string{}
}

func Interface(target interface{}, key string, typeStr []string) interface{} {

	if len(typeStr) != 0 && typeStr[0] == "string" {
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

func isSupportExt(s string) bool {
	return isJSONExt(s) || isYAMLExt(s)
}

func isJSONExt(s string) bool {
	ext := path.Ext(s)
	return ext == ".json"
}

func isYAMLExt(s string) bool {
	ext := path.Ext(s)
	return ext == ".yaml" || ext == ".yml"
}

func baseResourceName(s string) string {
	ss := strings.Split(s, "/")
	name := ss[len(ss)-1]

	ext := path.Ext(name)
	return name[0 : len(name)-len(ext)]
}

func isExtJSONFile(info os.FileInfo) bool {
	return isFile(info) && isMatchExt(info, ".json")
}

func isExtYaml(info os.FileInfo) bool {
	return isFile(info) && isMatchExt(info, ".yml", ".yaml")
}

func isMatchExt(info os.FileInfo, exts ...string) bool {
	e := filepath.Ext(info.Name())
	for _, ext := range exts {
		if e == ext {
			return true
		}
	}
	return false
}

func isFile(info os.FileInfo) bool {
	if info == nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func resourceName(name string) string {
	base := path.Base(name)
	e := filepath.Ext(base)
	return base[0 : len(base)-len(e)]
}

func YamlFileToJson(path string, info os.FileInfo) ([]byte, error) {
	isJSON := isExtJSONFile(info)
	isYAML := isExtYaml(info)

	if !isJSON && !isYAML {
		return nil, fmt.Errorf("%v is not support file format", info.Name())
	}

	rs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var d map[string]interface{}
	switch {
	case isJSON:
		if err := json.Unmarshal(rs, &d); err != nil {
			return nil, err
		}
	case isYAML:
		if err := yaml.Unmarshal(rs, &d); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%v is not support file format", info.Name())
	}

	return json.MarshalIndent(d, "", "  ")
}
