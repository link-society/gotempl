package decoder

import (
	"fmt"

	"github.com/joho/godotenv"
)

type EnvDecoder struct{}

func (e EnvDecoder) Format() string {
	return "env"
}

func (e EnvDecoder) Shortcut() string {
	return "e"
}

func (e EnvDecoder) Decode(input []byte) (Data, error) {
	env, err := godotenv.Unmarshal(string(input))
	if err != nil {
		return nil, fmt.Errorf("[env-decoder] %s", err)
	}

	data := map[string]any{}
	for key, val := range env {
		data[key] = val
	}

	return data, nil
}
