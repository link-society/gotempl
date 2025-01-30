package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	coreIO "io"
	"os"
	"path/filepath"
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

	err = io.ExecuteTemplate([]string{"--data-env", "./tests/data.env"})
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

	err = io.ExecuteTemplate([]string{"--data-json", "./tests/data.json"})
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

	err = io.ExecuteTemplate([]string{"--data-yaml", "./tests/data.yaml"})
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

	err = io.ExecuteTemplate([]string{"--data-toml", "./tests/data.toml"})
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

	err = io.ExecuteTemplate([]string{})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "HELLOHELLOHELLOHELLOHELLO"
	assert.Equal(t, result, expected)
}

func TestIsDir(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte(`{{ isDir "internal" }}`))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	err = io.ExecuteTemplate([]string{})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "true"
	assert.Equal(t, result, expected)
}

func TestReadDir(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte(`{{ readDir "internal" }}`))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	err = io.ExecuteTemplate([]string{})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "[decoder io options template]"
	assert.Equal(t, result, expected)
}

func TestReadFile(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	path := "tests/data.yaml"
	_, err := stdinBuf.Write([]byte(
		fmt.Sprintf("{{ readFile \"%s\" }}", path)),
	)
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	err = io.ExecuteTemplate([]string{})
	if err != nil {
		t.Error(err)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := string(bytes)
	assert.Equal(t, result, expected)
}

func TestWalkDir(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte(`{{ walkDir ".github" }}`))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	err = io.ExecuteTemplate([]string{})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := fmt.Sprintf(
		"[workflows %s %s]",
		filepath.Join("workflows", "release.yml"),
		filepath.Join("workflows", "test-suite.yml"),
	)
	assert.Equal(t, result, expected)
}

func TestFileExists(t *testing.T) {
	stdinBuf := new(bytes.Buffer)
	stdoutBuf := new(bytes.Buffer)

	_, err := stdinBuf.Write([]byte(`{{ fileExists "internal/decoder/env.go" }}`))
	if err != nil {
		t.Fatal(err)
	}

	io.SetInput(stdinBuf)
	io.SetOutput(stdoutBuf)

	err = io.ExecuteTemplate([]string{})
	if err != nil {
		t.Error(err)
	}

	result := stdoutBuf.String()
	expected := "true"
	assert.Equal(t, result, expected)
}

func TestTemplates(t *testing.T) {
	stdoutBuf := new(bytes.Buffer)

	io.SetOutput(stdoutBuf)

	err := io.ExecuteTemplate([]string{
		"-t", "tests/tmpl",
		"-t", "tests/tmpl",
	})
	if err != nil {
		t.Error(err)
	}

	reader := bufio.NewReader(stdoutBuf)
	line, ok, err := reader.ReadLine()
	assert.Equal(t, line, []byte("true"))
	assert.False(t, ok)
	assert.Nil(t, err)

	line, ok, err = reader.ReadLine()
	assert.Equal(t, line, []byte("true"))
	assert.False(t, ok)
	assert.Nil(t, err)

	line, ok, err = reader.ReadLine()
	assert.Nil(t, line)
	assert.False(t, ok)
	assert.ErrorIs(t, err, coreIO.EOF)
}
