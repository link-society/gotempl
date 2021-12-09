package internal

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hellflame/argparse"
)

type Options struct {
	Template   io.Reader
	Output     io.Writer
	DataParser *DataParser
}

func NewOptions(args []string) (opts Options, err error) {
	opts.Template = os.Stdin
	opts.Output = os.Stdout

	parser := argparse.NewParser(
		"gotempl", "Generic templating tool which use both environment variables and data files as template data",
		&argparse.ParserConfig{
			AddShellCompletion:     true,
			DisableDefaultShowHelp: true,
		},
	)

	parser.String(
		"", "template",
		&argparse.Option{
			Positional: true,
			Required:   false,
			Help:       "Path to Go Template file. Default is stdin. Caution: if you a template argument just after a data file argument, the template will be parsed as a data file. Example: \"TEST env var is {{ .Env.TEST }} and TEST data value is {{ .Data.TEST }}.\"",
			Validate: func(arg string) (err error) {
				opts.Template, err = os.Open(arg)

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
				opts.Output, err = os.Create(arg)

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
