package io

import (
	"io"
	"os"

	"github.com/link-society/gotempl/internal/decoder"
	"github.com/link-society/gotempl/internal/options"
	"github.com/link-society/gotempl/internal/template"
)

type Context struct {
	Template   *template.GenericTemplate
	Data       decoder.Data
	OutputPath string
}

func ExecuteTemplate(args []string) error {
	opts, err := options.ParseOptions(args)
	if err != nil {
		return err
	}

	context, err := NewContext(opts)
	if err != nil {
		return err
	}

	err = context.Write()
	if err != nil {
		return err
	}

	return nil
}

func NewContext(opts options.Options) (Context, error) {
	template, err := ReadTemplate(opts.TemplatePaths, opts.HTML)
	if err != nil {
		return Context{}, err
	}

	ctx := Context{
		Template: template,
		Data: map[string]interface{}{
			"Env":  ReadEnvironment(),
			"Data": opts.TemplateData,
		},
		OutputPath: opts.OutputPath,
	}

	return ctx, nil
}

func (ctx Context) Write() error {
	var writer io.Writer

	if ctx.OutputPath == "" {
		writer = Stdout()
	} else {
		file, err := os.Create(ctx.OutputPath)
		if err != nil {
			return err
		}

		writer = file

		defer file.Close()
	}

	return ctx.Template.Execute(writer, ctx.Data)
}
