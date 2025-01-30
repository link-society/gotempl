package io

import (
	"bytes"
	"fmt"
	"io"
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
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("[template-open] %s", err)
			}

			templateContent = append(templateContent, content...)
		}

		reader = bytes.NewReader(templateContent)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("[template-read] %s", err)
	}

	template, err := template.New("template", html).Parse(string(content))

	if err != nil {
		return nil, fmt.Errorf("[template-parse] %s", err)
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
