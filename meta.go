package main

import (
	"encoding/json"
	"io/ioutil"
)

type Meta struct {
	Title       string   `json:"title"`
	BaseURL     string   `json:"base_url"`
	ContentType string   `json:"content_type"`
	Headers     []string `json:"headers"`
}

func readMeta(path string) (Meta, error) {
	meta := Meta{
		Title:       "API Document",
		BaseURL:     "http://localhost",
		ContentType: "application/json",
		Headers: []string{
			"Content-Type: application/json",
		},
	}
	if path == "" {
		return meta, nil
	}
	rs, err := ioutil.ReadFile(path)

	if err != nil {
		return meta, err
	}

	if err := json.Unmarshal(rs, &meta); err != nil {
		return meta, err
	}

	return meta, nil
}
