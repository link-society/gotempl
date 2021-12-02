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

func ReadInputFiles(opts Options) (ctx Context, err error) {
	templateContent, err := ioutil.ReadAll(opts.Template)
	if err != nil {
		return
	}

	ctx.Template, err = template.New("template").Parse(string(templateContent))
	if err != nil {
		return
	}

	data, err := opts.DataParser.GetData()
	if err != nil {
		return
	}

	env := getEnvironment()

	ctx.Data = map[string]interface{}{
		"Data": data,
		"Env":  env,
	}

	return ctx, err
}

func getEnvironment() map[string]string {
	var data = map[string]string{}

	for _, env := range os.Environ() {
		var splitted = strings.SplitN(env, "=", 2)
		data[splitted[0]] = splitted[1]
	}

	return data
}

func WriteOutput(opts Options, context Context) (err error) {
	return context.Template.Execute(opts.Output, context.Data)
}
