package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/akamensky/argparse"
	"github.com/go-yaml/yaml"
)

type Options struct {
	templatePath string
	dataPath     string
	dataFormat   string
	outputPath   string
}

type Context struct {
	template *template.Template
	data     map[string]interface{}
}

func parseArgs(args []string) (*Options, error) {
	parser := argparse.NewParser("gotempl", "Generic templating tool")

	templatePath := parser.String("t", "template", &argparse.Options{
		Required: true,
		Help:     "Path to Go Template file",
	})
	dataPath := parser.String("d", "data", &argparse.Options{
		Required: true,
		Help:     "Path to data file to use for templating",
	})
	dataFormat := parser.String("f", "format", &argparse.Options{
		Required: false,
		Help:     "Format of data file (json, yaml or toml, defaults to json)",
		Default:  "json",
	})
	outputPath := parser.String("o", "output", &argparse.Options{
		Required: false,
		Help:     "Path to output file (leave empty for stdout)",
		Default:  "",
	})

	err := parser.Parse(args)

	if err != nil {
		return nil, errors.New(parser.Usage(err))
	}

	opts := &Options{
		templatePath: *templatePath,
		dataPath:     *dataPath,
		dataFormat:   *dataFormat,
		outputPath:   *outputPath,
	}

	return opts, nil
}

func readInputFiles(opts *Options) (*Context, error) {
	templateContent, err := ioutil.ReadFile(opts.templatePath)
	if err != nil {
		return nil, err
	}

	dataContent, err := ioutil.ReadFile(opts.dataPath)
	if err != nil {
		return nil, err
	}

	template, err := template.New(opts.templatePath).Parse(string(templateContent))
	if err != nil {
		return nil, err
	}

	data, err := decodeInputData(dataContent, opts.dataFormat)
	if err != nil {
		return nil, err
	}

	context := &Context{
		template: template,
		data:     data,
	}

	return context, err
}

func decodeInputData(buf []byte, format string) (map[string]interface{}, error) {
	var (
		result map[string]interface{}
		err    error
	)

	switch format {
	case "json":
		err = json.Unmarshal(buf, &result)

	case "yaml":
		err = yaml.Unmarshal(buf, &result)

	case "toml":
		err = toml.Unmarshal(buf, &result)

	default:
		err = errors.New("Unsupported file type")
	}

	return result, err
}

func writeOutput(path string, context *Context) error {
	var writer io.Writer

	if path == "" {
		writer = os.Stdout
	} else {
		f, err := os.Create(path)
		if err != nil {
			return err
		}

		defer f.Close()
		writer = f
	}

	return context.template.Execute(writer, context.data)
}

func main() {
	opts, err := parseArgs(os.Args)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	context, err := readInputFiles(opts)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	err = writeOutput(opts.outputPath, context)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
