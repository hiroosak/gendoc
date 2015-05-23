# Gendoc

[![Build Status](https://travis-ci.org/hiroosak/gendoc.svg?branch=master)](https://travis-ci.org/hiroosak/gendoc)

Gendoc is generate documentation from [JSON Schema](http://json-schema.org/)

This is inspired by [Prmd](https://github.com/interagent/prmd).

## Install

``` bash
go get github.com/hiroosak/gendoc
```

## Usage

Gendoc provides these commands.

* `init` - Create initialized YAML file
* `doc` - Generate HTML from json schema
* `valid` - Validation JSON Schema format 
* `gen` - Generate JSON from YAML

### Example

``` bash
$ gendoc init article > src/article.yml
$ gendoc init comment > src/comment.yml
$ vim src/{article,comment}.yaml

# Build docs
$ gendoc doc -src ./src -meta meta.json > docs.html
```

## doc

Generate a HTML document. This command has these flag options.

* `src` - directory where the yaml, json file entered
* `meta` - overall API metadata
* `overview` - preamble for generated API docs(html format)

``` bash
# Build docs
$ gendoc doc -src ./src -meta meta.json -overview overview.html > docs.html
```

### meta 

``` json
{
  "title": "API Title",
  "base_url": "http://localhost/",
  "content_type": "application/json",
  "headers": [
    "X-Service-Token: AAA"
  ]
}
```

## YAML to JSON

Convert the yaml files under the src directory to JSON.

``` bash
# Build json
$ gendoc gen -src ./src -dst ./dst # src/article.yml -> dst/articles.json
```

# License

MIT
