package main

import (
	"bytes"
	"os"
	"testing"
)

const expectedOutput = "foo is bar and prop is val while test is true"

func runTest(t *testing.T, opts Options) {
	os.Setenv("TEST", "true")

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

	if s != expectedOutput {
		t.Errorf("Template generation failed: %v. Expected: %v", s, expectedOutput)
	}
}

func TestJSON(t *testing.T) {
	runTest(t, Options{
		templatePath: "./tests/example.tmpl",
		dataPath:     "./tests/example.json",
		dataFormat:   "json",
		outputPath:   "",
	})
}

func TestYAML(t *testing.T) {
	runTest(t, Options{
		templatePath: "./tests/example.tmpl",
		dataPath:     "./tests/example.yml",
		dataFormat:   "yaml",
		outputPath:   "",
	})
}

func TestTOML(t *testing.T) {
	runTest(t, Options{
		templatePath: "./tests/example.tmpl",
		dataPath:     "./tests/example.toml",
		dataFormat:   "toml",
		outputPath:   "",
	})
}

func TestENV(t *testing.T) {
	runTest(t, Options{
		templatePath: "./tests/example.env.tmpl",
		dataPath:     "./tests/example.env",
		dataFormat:   "env",
		outputPath:   "",
	})
}
