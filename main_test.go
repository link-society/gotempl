package main_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/link-society/gotempl/internal"
)

func NewTestData(prefixTemplate, prefixExpected string) (dataParser internal.DataParser, template string, expected string) {
	template = prefixTemplate
	expected = prefixExpected

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

func Test(t *testing.T) {
	os.Setenv("TEST", "true")

	var opts = internal.Options{}
	dataParser, template, expected := NewTestData(template, expected)
	opts.DataParser = dataParser

	const templatePath = "/tmp/gotempl.test"
	templateFile, err := os.Create(templatePath)

	if err != nil {
		t.Fatal(err)
	}

	_, err = templateFile.WriteString(template)

	if err != nil {
		t.Fatal(err)
	}

	templateFile.Seek(0, 0)

	opts.Template = templateFile

	context, err := opts.ReadInputFiles()

	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(templatePath)

	if err != nil {
		t.Fatal(err)
	}

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
