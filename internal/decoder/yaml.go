package decoder

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type YamlDecoder struct{}

func (y YamlDecoder) Format() string {
	return "yaml"
}

func (y YamlDecoder) Shortcut() string {
	return "y"
}

func (y YamlDecoder) Decode(input []byte) (Data, error) {
	var data Data

	err := yaml.Unmarshal(input, &data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[yaml-decoder] %s", err))
	}

	return data, nil
}
