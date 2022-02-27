package template

import (
	"io"

	htemplate "html/template"
	ttemplate "text/template"

	"github.com/Masterminds/sprig"
)

type GenericFuncMap = map[string]interface{}

type GenericTemplate struct {
	isHTML bool

	htmlTemplate *htemplate.Template
	textTemplate *ttemplate.Template
}

func New(name string, html bool) *GenericTemplate {
	if html {
		return NewHTML(name)
	} else {
		return NewText(name)
	}
}

func NewHTML(name string) *GenericTemplate {
	return &GenericTemplate{
		isHTML:       true,
		htmlTemplate: htemplate.New(name).Funcs(sprig.HtmlFuncMap()).Funcs(funcs),
		textTemplate: nil,
	}
}

func NewText(name string) *GenericTemplate {
	return &GenericTemplate{
		isHTML:       false,
		htmlTemplate: nil,
		textTemplate: ttemplate.New(name).Funcs(sprig.TxtFuncMap()).Funcs(funcs),
	}
}

func (tmpl *GenericTemplate) Parse(text string) (*GenericTemplate, error) {
	var err error

	if tmpl.isHTML {
		tmpl.htmlTemplate, err = tmpl.htmlTemplate.Parse(text)
	} else {
		tmpl.textTemplate, err = tmpl.textTemplate.Parse(text)
	}

	return tmpl, err
}

func (tmpl *GenericTemplate) Execute(wr io.Writer, data interface{}) error {
	var err error

	if tmpl.isHTML {
		err = tmpl.htmlTemplate.Execute(wr, data)
	} else {
		err = tmpl.textTemplate.Execute(wr, data)
	}

	return err
}
