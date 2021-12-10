package io

import (
	"html/template"
	"io"
	"os"

	"github.com/link-society/gotempl/internal/decoder"
	"github.com/link-society/gotempl/internal/options"
)

type Context struct {
	Template   *template.Template
	Data       decoder.Data
	OutputPath string
}

func NewContext(opts options.Options) (Context, error) {
	template, err := ReadTemplate(opts.TemplatePath)
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
