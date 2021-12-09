package decoder

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

type TomlDecoder struct{}

func (t TomlDecoder) Format() string {
	return "toml"
}

func (t TomlDecoder) Shortcut() string {
	return "T"
}

func (t TomlDecoder) Decode(input []byte) (Data, error) {
	var data Data

	err := toml.Unmarshal(input, &data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[toml-decoder] %s", err))
	}

	return data, nil
}
