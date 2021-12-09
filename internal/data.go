package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hellflame/argparse"
)

type Data map[string]interface{}

type DataParser struct {
	ArgDataParsers []ArgDataParser
}

func (dataParser DataParser) String() string {
	return fmt.Sprintf("Parsers: %v", dataParser.ArgDataParsers)
}

func NewDataParser(argParser *argparse.Parser) *DataParser {
	dataParser := &DataParser{}

	for format, decoder := range DecodersByFormat {
		dataParser.NewArgDataParser(format, decoder, argParser)
	}

	return dataParser
}

func (dataParser *DataParser) GetData() (Data, error) {
	data := Data{}

	for _, dataParser := range dataParser.ArgDataParsers {
		reader := dataParser.GetNextReader()
		bytes, err := ioutil.ReadAll(reader)

		if err != nil {
			return nil, err
		}

		err = dataParser.Decoder(bytes, data)

		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func GetArgName(format string) string {
	return fmt.Sprintf("data-%v", format)
}

func GetShortcut(format string) string {
	return string(format[0])
}

type ArgDataParser struct {
	Readers []io.Reader
	Decoder DataDecoder
	index   int
}

func (argDataParser ArgDataParser) String() string {
	return fmt.Sprintf("Readers: %v, index: %v", argDataParser.Readers, argDataParser.index)
}

func (argDataParser *ArgDataParser) GetNextReader() io.Reader {
	reader := argDataParser.Readers[argDataParser.index]
	argDataParser.index += 1
	return reader
}

func (dataParser *DataParser) NewArgDataParser(format string, decoder DataDecoder, argParser *argparse.Parser) ArgDataParser {
	var arg = GetArgName(format)
	var shortcut = GetShortcut(format)
	var help = fmt.Sprintf("Path to %v data file to use for templating", format)

	argDataParser := ArgDataParser{
		Decoder: decoder,
	}

	argParser.Strings(
		shortcut,
		arg,
		&argparse.Option{
			Required: false,
			Help:     help,
			Validate: func(arg string) error {
				file, err := os.Open(arg)
				if err != nil {
					return err
				}

				argDataParser.Readers = append(argDataParser.Readers, file)
				dataParser.ArgDataParsers = append(dataParser.ArgDataParsers, argDataParser)

				return nil
			},
		},
	)

	return argDataParser
}
