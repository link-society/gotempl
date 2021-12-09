package decoder

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JsonDecoder struct{}

func (j JsonDecoder) Format() string {
	return "json"
}

func (j JsonDecoder) Shortcut() string {
	return "j"
}

func (j JsonDecoder) Decode(input []byte) (Data, error) {
	var data Data

	err := json.Unmarshal(input, &data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[json-decoder] %s", err))
	}

	return data, nil
}
