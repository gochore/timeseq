package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"
	"text/template"
)

type Meta struct {
	Name string
	Type string
}

func main() {
	tmpl, err := template.ParseFiles("cmd/generate/x_sequence.go.tmpl")
	if err != nil {
		panic(err)
	}

	types := []string{
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int8",
		"int16",
		"int32",
		"int64",
		"float32",
		"float64",
		"int",
		"uint",
		"uint8",
		"int32",
	}

	for _, v := range types {
		generate( tmpl, Meta{
			Name: strings.Title(v),
			Type: v,
		})
	}
}

func generate(tmpl *template.Template, meta Meta) {
	output := &bytes.Buffer{}
	if err := tmpl.Execute(output, meta); err != nil {
		panic(err)
	}

	src, err := format.Source(output.Bytes())
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s_sequence.go", meta.Type), src, 0666); err != nil {
		panic(err)
	}
}