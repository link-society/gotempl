package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/link-society/gotempl/internal/io"
	"github.com/stretchr/testify/assert"
)

func TestEnvDataFile(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte("{{ .Env.PREFIX }}: {{ .Data.format }} is {{ .Data.env }}"))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	os.Setenv("PREFIX", "TESTENV")

	err = ExecuteTemplate([]string{"--data-env", "./tests/data.env"})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "TESTENV: env is test env"
	assert.Equal(t, result, expected)
}

func TestJsonDataFile(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte("{{ .Env.PREFIX }}: {{ .Data.format }} is {{ .Data.json }}"))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	os.Setenv("PREFIX", "TESTJSON")

	err = ExecuteTemplate([]string{"--data-json", "./tests/data.json"})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "TESTJSON: json is test json"
	assert.Equal(t, result, expected)
}

func TestYamlDataFile(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte("{{ .Env.PREFIX }}: {{ .Data.format }} is {{ .Data.yaml }}"))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	os.Setenv("PREFIX", "TESTYAML")

	err = ExecuteTemplate([]string{"--data-yaml", "./tests/data.yaml"})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "TESTYAML: yaml is test yaml"
	assert.Equal(t, result, expected)
}

func TestTomlDataFile(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte("{{ .Env.PREFIX }}: {{ .Data.format }} is {{ .Data.toml }}"))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	os.Setenv("PREFIX", "TESTTOML")

	err = ExecuteTemplate([]string{"--data-toml", "./tests/data.toml"})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "TESTTOML: toml is test toml"
	assert.Equal(t, result, expected)
}

func TestSprig(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte("{{ .Env.VALUE | upper | repeat 5 }}"))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	os.Setenv("VALUE", "hello")

	err = ExecuteTemplate([]string{})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "HELLOHELLOHELLOHELLOHELLO"
	assert.Equal(t, result, expected)
}
