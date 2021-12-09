package internal

import (
	"fmt"
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

func NewDataParser(argParser *argparse.Parser) (dataParser *DataParser) {
	dataParser = &DataParser{}

	for format, decoder := range DecodersByFormat {
		dataParser.NewArgDataParser(format, decoder, argParser)
	}

	return
}

func (dataParser *DataParser) GetData() (data Data, err error) {
	data = Data{}
	var bytes []byte
	var file *os.File

	for _, dataParser := range dataParser.ArgDataParsers {
		file, err = dataParser.GetNextFile()

		if err != nil {
			return
		}

		bytes, err = ioutil.ReadAll(file)

		if err != nil {
			return
		}

		err = dataParser.Decoder(bytes, data)

		if err != nil {
			return
		}
	}

	return
}

func GetArgName(format string) string {
	return fmt.Sprintf("data-%v", format)
}

func GetShortcut(format string) string {
	return string(format[0])
}

type ArgDataParser struct {
	Files   []*os.File
	Decoder DataDecoder
	index   int
}

func (argDataParser ArgDataParser) String() string {
	return fmt.Sprintf("Files: %v, index: %v", argDataParser.Files, argDataParser.index)
}

func (argDataParser *ArgDataParser) GetNextFile() (file *os.File, err error) {
	file = argDataParser.Files[argDataParser.index]
	argDataParser.index += 1
	return
}

func (dataParser *DataParser) NewArgDataParser(format string, decoder DataDecoder, argParser *argparse.Parser) (argDataParser ArgDataParser) {
	var arg = GetArgName(format)
	var shortcut = GetShortcut(format)
	var help = fmt.Sprintf("Path to %v data file to use for templating", format)

	argDataParser = ArgDataParser{
		Decoder: decoder,
	}

	argParser.Strings(
		shortcut,
		arg,
		&argparse.Option{
			Required: false,
			Help:     help,
			Validate: func(arg string) (err error) {
				var file *os.File
				file, err = os.Open(arg)

				if err == nil {
					argDataParser.Files = append(argDataParser.Files, file)
					dataParser.ArgDataParsers = append(dataParser.ArgDataParsers, argDataParser)
				}

				return
			},
		},
	)

	return
}
