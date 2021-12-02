package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/hellflame/argparse"
)

type Options struct {
	Template   *os.File
	Output     *os.File
	DataParser DataParser
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
