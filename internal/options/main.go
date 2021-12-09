package options

import (
	"errors"
	"fmt"

	"github.com/hellflame/argparse"
	"github.com/imdario/mergo"
	"github.com/link-society/gotempl/internal/decoder"
)

type Options struct {
	TemplatePath string
	TemplateData decoder.Data
	OutputPath   string
}

func ParseOptions(args []string) (Options, error) {
	opts := Options{
		TemplatePath: "",
		TemplateData: map[string]interface{}{},
		OutputPath:   "",
	}

	parser := argparse.NewParser(
		"gotempl", "Generic templating tool which use both environment variables and data files as template data",
		&argparse.ParserConfig{
			AddShellCompletion:     true,
			DisableDefaultShowHelp: true,
		},
	)

	parser.String(
		"t", "template",
		&argparse.Option{
			Positional: false,
			Required:   false,
			Help:       "Path to Go Template file. Default is stdin.",
			Validate: func(arg string) error {
				opts.TemplatePath = arg
				return nil
			},
		},
	)

	parser.String(
		"o", "output",
		&argparse.Option{
			Required: false,
			Help:     "Path to output file. Default is stdout",
			Validate: func(arg string) error {
				opts.OutputPath = arg
				return nil
			},
		},
	)

	decoder.AddOptions(parser, func(data decoder.Data) error {
		err := mergo.Merge(&opts.TemplateData, data)
		if err != nil {
			return errors.New(fmt.Sprintf("[data-merge] %s", err))
		}

		return nil
	})

	err := parser.Parse(args)
	if err != nil {
		return Options{}, errors.New(
			fmt.Sprintf("%v\n\n%v", err, parser.FormatHelp()),
		)
	}

	return opts, nil
}
