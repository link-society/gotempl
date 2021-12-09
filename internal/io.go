package internal

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type Context struct {
	Template *template.Template
	Data     map[string]interface{}
}

func ReadInputFiles(opts *Options) (*Context, error) {
	templateContent, err := ioutil.ReadAll(opts.Template)
	if err != nil {
		return nil, err
	}

	template, err := template.New("template").Parse(string(templateContent))
	if err != nil {
		return nil, err
	}

	data, err := opts.DataParser.GetData()
	if err != nil {
		return nil, err
	}

	env := getEnvironment()

	ctx := &Context{
		Template: template,
		Data: map[string]interface{}{
			"Data": data,
			"Env":  env,
		},
	}

	return ctx, nil
}

func getEnvironment() map[string]string {
	var data = map[string]string{}

	for _, env := range os.Environ() {
		var splitted = strings.SplitN(env, "=", 2)
		data[splitted[0]] = splitted[1]
	}

	return data
}

func WriteOutput(opts *Options, context *Context) error {
	return context.Template.Execute(opts.Output, context.Data)
}
