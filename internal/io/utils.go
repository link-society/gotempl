package io

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/link-society/gotempl/internal/template"
)

func ReadTemplate(paths []string, html bool) (*template.GenericTemplate, error) {
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

	template, err := template.New("template", html).Parse(string(content))

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
