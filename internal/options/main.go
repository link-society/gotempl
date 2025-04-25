package options

import (
	"fmt"
	"strings"

	"dario.cat/mergo"
	"github.com/google/shlex"
	"github.com/hellflame/argparse"
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
		TemplateData: map[string]any{},
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

	parser.Strings(
		"d", "data",
		&argparse.Option{
			Positional: false,
			Required:   false,
			Help:       "Data value of format name=value to be used in template. Can be used multiple times.",
			Validate: func(arg string) error {
				toks, err := shlex.Split(arg)
				if err != nil {
					return fmt.Errorf("Error: %s on data option: %s.", err.Error(), arg)
				}
				if len(toks) == 0 {
					return fmt.Errorf("data-value must be in the form of key=value. Found %s", arg)
				}

				keyAndValue := strings.SplitN(toks[0], "=", 2)
				if len(keyAndValue) != 2 {
					return fmt.Errorf("data-value must be in the form of key=value. Found %s", arg)
				}
				opts.TemplateData[strings.TrimSpace(keyAndValue[0])] = strings.TrimSpace(keyAndValue[1])

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
