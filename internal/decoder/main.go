package decoder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hellflame/argparse"
)

type Data map[string]interface{}
type DataContinuation func(Data) error

type Decoder interface {
	Format() string
	Shortcut() string
	Decode(input []byte) (Data, error)
}

var decoders = []Decoder{
	JsonDecoder{},
	YamlDecoder{},
	TomlDecoder{},
	EnvDecoder{},
}

func AddOptions(parser *argparse.Parser, cont DataContinuation) {
	for _, decoder := range decoders {
		parser.Strings(
			decoder.Shortcut(),
			fmt.Sprintf("data-%s", decoder.Format()),
			&argparse.Option{
				Required: false,
				Help:     fmt.Sprintf("Path to %s file", strings.ToUpper(decoder.Format())),
				Validate: NewValidator(decoder, cont),
			},
		)
	}
}

func NewValidator(decoder Decoder, cont DataContinuation) func(string) error {
	return func(arg string) error {
		file, err := os.Open(arg)
		if err != nil {
			return errors.New(fmt.Sprintf("[decoder-open] %s", err))
		}

		defer file.Close()

		buf, err := ioutil.ReadAll(file)
		if err != nil {
			return errors.New(fmt.Sprintf("[decoder-read] %s", err))
		}

		data, err := decoder.Decode(buf)
		if err != nil {
			return err
		}

		err = cont(data)
		if err != nil {
			return err
		}

		return nil
	}
}
