package io

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Masterminds/sprig"
)

func ReadTemplate(path string) (*template.Template, error) {
	var reader io.Reader

	if path == "" {
		reader = Stdin()
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("[template-open] %s", err))
		}

		reader = file

		defer file.Close()
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[template-read] %s", err))
	}

	template, err := template.
		New("template").
		Funcs(sprig.FuncMap()).
		Parse(string(content))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[template-parse] %s", err))
	}

	return template, nil
}

func ReadEnvironment() map[string]string {
	var data = map[string]string{}

	for _, env := range os.Environ() {
		var splitted = strings.SplitN(env, "=", 2)
		data[splitted[0]] = splitted[1]
	}

	return data
}
