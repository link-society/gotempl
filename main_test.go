package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/link-society/gotempl/internal"
)

func NewTestCLI(prefixTemplate, prefixExpected string) (args []string, template string, expected string) {
	template = prefixTemplate
	expected = prefixExpected

	var format string

	for format, _ = range internal.DecodersByFormat {
		var filename = fmt.Sprintf("./tests/data.%v", format)
		var arg = fmt.Sprintf("--data-%v", format)
		args = append(args, arg, filename)

		template = fmt.Sprintf("%v, {{ .Data.%v }}", template, format)
		expected = fmt.Sprintf("%v, test %v", expected, format)
	}

	expected = fmt.Sprintf("%v %v", format, expected)

	return
}

func NewTestData(prefixTemplate, prefixExpected string) (dataParser *internal.DataParser, template string, expected string) {
	template = prefixTemplate
	expected = prefixExpected
	dataParser = &internal.DataParser{}

	var format string
	var decoder internal.DataDecoder
	for format, decoder = range internal.DecodersByFormat {
		var filename = fmt.Sprintf("./tests/data.%v", format)
		var file, err = os.Open(filename)
		if err != nil {
			panic(err)
		}
		var argDataParser = internal.ArgDataParser{
			Files:   []*os.File{file},
			Decoder: decoder,
		}

		dataParser.ArgDataParsers = append(dataParser.ArgDataParsers, argDataParser)

		template = fmt.Sprintf("%v, {{ .Data.%v }}", template, format)
		expected = fmt.Sprintf("%v, test %v", expected, format)
	}

	expected = fmt.Sprintf("%v %v", format, expected)

	return
}

const template = "{{ .Data.format }} {{ .Env.TEST }}"

const expected = "true"

func NewTemplate(t *testing.T, content string) (template *os.File) {
	const templatePath = "gotempl.test"
	template, err := os.CreateTemp("tests", templatePath)

	if err != nil {
		t.Fatal(err)
	}

	_, err = template.WriteString(content)

	if err != nil {
		t.Fatal(err)
	}

	_, err = template.Seek(0, 0)

	if err != nil {
		t.Fatal(err)
	}

	return
}

func TestCLI(t *testing.T) {
	os.Setenv("TEST", "true")

	args, template, expected := NewTestCLI(template, expected)

	var templateFile = NewTemplate(t, template)

	var cmdArgs = []string{
		"run",
		"main.go",
		templateFile.Name(),
	}

	cmdArgs = append(cmdArgs, args...)

	var cmd = exec.Command("go", cmdArgs...)

	var stdout, _ = cmd.Output()
	var output = string(stdout)

	if output != expected {
		t.Fatalf("Template generation failed: %v. Expected: %v", string(output), expected)
	}

	os.Remove(templateFile.Name())
}

func TestCLIIN(t *testing.T) {
	os.Setenv("TEST", "true")

	args, template, expected := NewTestCLI(template, expected)

	var cmdArgs = []string{
		"run",
		"main.go",
	}

	cmdArgs = append(cmdArgs, args...)

	var cmd = exec.Command("go", cmdArgs...)

	var stdin, _ = cmd.StdinPipe()

	var _, err = stdin.Write([]byte(template))

	if err != nil {
		t.Fatal(err)
	}

	err = stdin.Close()

	if err != nil {
		t.Fatal(err)
	}

	var stdout, _ = cmd.Output()
	var output = string(stdout)

	if output != expected {
		t.Fatalf("Template generation failed: %v. Expected: %v", string(output), expected)
	}
}

func TestCLIOUT(t *testing.T) {
	os.Setenv("TEST", "true")

	args, template, expected := NewTestCLI(template, expected)

	var outputFile = "output.result"

	var cmdArgs = []string{
		"run",
		"main.go",
		"-o",
		outputFile,
	}

	cmdArgs = append(cmdArgs, args...)

	var cmd = exec.Command("go", cmdArgs...)

	var stdin, _ = cmd.StdinPipe()

	var _, err = stdin.Write([]byte(template))

	if err != nil {
		t.Fatal(err)
	}

	err = stdin.Close()

	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Run()

	if err != nil {
		t.Fatal(err)
	}

	var outputBytes []byte
	outputBytes, err = ioutil.ReadFile(outputFile)

	if err != nil {
		t.Fatal(err)
	}

	os.Remove(outputFile)

	var output = string(outputBytes)

	if output != expected {
		t.Fatalf("Template generation failed: %v. Expected: %v", output, expected)
	}

}

func TestData(t *testing.T) {
	os.Setenv("TEST", "true")

	var opts = internal.Options{}
	dataParser, template, expected := NewTestData(template, expected)
	opts.DataParser = dataParser

	var templateFile = NewTemplate(t, template)

	opts.Template = templateFile

	context, err := internal.ReadInputFiles(opts)

	if err != nil {
		t.Fatal(err)
	}

	os.Remove(templateFile.Name())

	buf := new(bytes.Buffer)
	err = context.Template.Execute(buf, context.Data)

	if err != nil {
		t.Fatal(err)
	}

	s := buf.String()

	if s != expected {
		t.Fatalf("Template generation failed: %v. Expected: %v", s, expected)
	}
}
