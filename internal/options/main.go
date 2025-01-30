package options

import (
	"fmt"

	"github.com/hellflame/argparse"
	"github.com/imdario/mergo"
	"github.com/link-society/gotempl/internal/decoder"
)

type Options struct {
	TemplatePaths []string
	TemplateData  decoder.Data
	OutputPath    string
	HTML          bool
}

func ParseOptions(args []string) (Options, error) {
	opts := Options{
		TemplateData: map[string]interface{}{},
	}

	parser := argparse.NewParser(
		"gotempl", "Generic templating tool which use both environment variables and data files as template data",
		&argparse.ParserConfig{
			AddShellCompletion:     true,
			DisableDefaultShowHelp: true,
		},
	)

	htmlFlag := parser.Flag(
		"H", "html",
		&argparse.Option{
			Positional: false,
			Required:   false,
			Help:       "Escape template for HTML output",
		},
	)

	parser.Strings(
		"t", "template",
		&argparse.Option{
			Positional: false,
			Required:   false,
			Help:       "Path to Go Template file. Default is stdin.",
			Validate: func(arg string) error {
				opts.TemplatePaths = append(opts.TemplatePaths, arg)
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
			return fmt.Errorf("[data-merge] %s", err)
		}

		return nil
	})

	err := parser.Parse(args)
	if err != nil {
		return Options{}, fmt.Errorf("%v\n\n%v", err, parser.FormatHelp())
	}

	opts.HTML = *htmlFlag

	return opts, nil
}
