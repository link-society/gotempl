package io

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Masterminds/sprig"
)

func ReadTemplate(paths []string) (*template.Template, error) {
	var reader io.Reader

	if len(paths) == 0 {
		reader = Stdin()
	} else {
		var templateContent []byte

		for _, path := range paths {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("[template-open] %s", err))
			}

			templateContent = append(templateContent, content...)
		}

		reader = bytes.NewReader(templateContent)
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[template-read] %s", err))
	}

	template, err := template.
		New("template").
		Funcs(sprig.FuncMap()).
		Funcs(Funcs).
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
