package main

import (
	"bytes"
	"testing"
)

func runTest(t *testing.T, opts *Options) {
	context, err := readInputFiles(opts)

	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)
	err = context.template.Execute(buf, context.data)
	if err != nil {
		t.Error(err)
	}

	s := buf.String()

	if s != "foo is bar and prop is val" {
		t.Errorf("Template generation failed")
	}
}

func TestJSON(t *testing.T) {
	runTest(t, &Options{
		templatePath: "./tests/example.tmpl",
		dataPath:     "./tests/example.json",
		dataFormat:   "json",
		outputPath:   "",
	})
}

func TestYAML(t *testing.T) {
	runTest(t, &Options{
		templatePath: "./tests/example.tmpl",
		dataPath:     "./tests/example.yml",
		dataFormat:   "yaml",
		outputPath:   "",
	})
}

func TestTOML(t *testing.T) {
	runTest(t, &Options{
		templatePath: "./tests/example.tmpl",
		dataPath:     "./tests/example.toml",
		dataFormat:   "toml",
		outputPath:   "",
	})
}
