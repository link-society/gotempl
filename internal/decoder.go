package internal

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/joho/godotenv"
)

// DataDecoder fill input map data with input bytes
type DataDecoder = func(input []byte, data Data) error

var DecodersByFormat = map[string]DataDecoder{
	"json": jsonDecoder,
	"yaml": yamlDecoder,
	"toml": tomlDecoder,
	"env":  envDecoder,
}

func jsonDecoder(input []byte, data Data) error {
	return json.Unmarshal(input, &data)
}

func yamlDecoder(input []byte, data Data) error {
	return yaml.Unmarshal(input, &data)
}

func tomlDecoder(input []byte, data Data) error {
	return toml.Unmarshal(input, &data)
}

func envDecoder(input []byte, data Data) error {
	envMap, err := godotenv.Unmarshal(string(input))

	if err != nil {
		return err
	}

	for key, value := range envMap {
		data[key] = value
	}

	return nil
}
