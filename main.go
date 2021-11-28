package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/akamensky/argparse"
	"github.com/go-yaml/yaml"
	"github.com/joho/godotenv"
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

func parseArgs(args []string) (Options, error) {
	var opts Options
	parser := argparse.NewParser("gotempl", "Generic templating tool")

	templatePath := parser.String("t", "template", &argparse.Options{
		Required: true,
		Help:     "Path to Go Template file.\n\t\t\tExample: \"TEST environment variable is {{ .Env.TEST }} and TEST data value is {{ .Data.TEST }}.\"",
	})
	dataPath := parser.String("d", "data", &argparse.Options{
		Required: true,
		Help:     "Path to data file to use for templating",
	})
	dataFormat := parser.String("f", "format", &argparse.Options{
		Required: false,
		Help:     "Format of data file (json, yaml, toml or env, defaults to json)",
		Default:  "json",
	})
	outputPath := parser.String("o", "output", &argparse.Options{
		Required: false,
		Help:     "Path to output file (leave empty for stdout)",
		Default:  "",
	})

	err := parser.Parse(args)

	if err != nil {
		return opts, errors.New(parser.Usage(err))
	}

	opts = Options{
		templatePath: *templatePath,
		dataPath:     *dataPath,
		dataFormat:   *dataFormat,
		outputPath:   *outputPath,
	}

	return opts, nil
}

func readInputFiles(opts Options) (*Context, error) {
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

	env := getEnvironment()

	context := &Context{
		template: template,
		data: map[string]interface{}{
			"Data": data,
			"Env":  env,
		},
	}

	return context, err
}

func getEnvironment() map[string]string {
	var data = map[string]string{}

	for _, env := range os.Environ() {
		var splitted = strings.SplitN(env, "=", 2)
		data[splitted[0]] = splitted[1]
	}

	return data
}

func unmarshalEnv(data []byte, v interface{}) error {
	envs, err := godotenv.Unmarshal(string(data))

	if err == nil {
		var vValue = reflect.ValueOf(v)
		var vMap = reflect.MakeMap(vValue.Elem().Type())

		if err == nil {
			for key, value := range envs {
				vMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
			}
		}

		vValue.Elem().Set(vMap)
	}

	return err
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

	case "env":
		err = unmarshalEnv(buf, &result)

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
