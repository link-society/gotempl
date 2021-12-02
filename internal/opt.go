package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/hellflame/argparse"
)

type Options struct {
	Template   *os.File
	Output     *os.File
	DataParser DataParser
}

type Context struct {
	Template *template.Template
	Data     map[string]interface{}
}

func NewOptions(args []string) (opts Options, err error) {
	parser := argparse.NewParser(
		"gotempl", "Generic templating tool",
		&argparse.ParserConfig{
			AddShellCompletion: true,
		},
	)

	parser.String(
		"", "template",
		&argparse.Option{
			Positional: true,
			Required:   false,
			Help:       "Path to Go Template file. Default is stdin. Example: \"TEST env var is {{ .Env.TEST }} and TEST data value is {{ .Data.TEST }}.\"",
			Validate: func(arg string) (err error) {
				if arg == "" {
					opts.Template = os.Stdin
				} else {
					opts.Template, err = os.Open(arg)
				}

				return
			},
		},
	)

	parser.String(
		"o", "output",
		&argparse.Option{
			Required: false,
			Help:     "Path to output file. Default is stdout",
			Validate: func(arg string) (err error) {
				if arg == "" {
					opts.Output = os.Stdout
				} else {
					opts.Output, err = os.Create(arg)
				}

				return
			},
		},
	)

	opts.DataParser = NewDataParser(parser)

	err = parser.Parse(args)

	if err != nil {
		err = errors.New(
			fmt.Sprintf("%v\n\n%v", err, parser.FormatHelp()),
		)
	}

	return
}

func (opts Options) ReadInputFiles() (ctx Context, err error) {
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

func (opts Options) WriteOutput(context Context) (err error) {
	return context.Template.Execute(opts.Output, context.Data)
}
