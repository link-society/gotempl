# gotempl

Small binary used to generate files from Go Templates, environment variables and data files.

The following formats are supported:

- JSON
- YAML
- TOML
- ENV (environment variables in the form "key=value")

## Usage

```bash
usage: gotempl [--help] [--completion] [--output OUTPUT] [--data-yaml DATA-YAML [DATA-YAML ...]] [--data-toml DATA-TOML [DATA-TOML ...]] [--data-env DATA-ENV [DATA-ENV ...]] [--data-json DATA-JSON [DATA-JSON ...]] [TEMPLATE]

Generic templating tool

positional arguments:
  TEMPLATE                             Path to Go Template file. Default is stdin. Example: "TEST env var is {{ .Env.TEST }} and TEST data value is {{ .Data.TEST
                                        }}."

optional arguments:
  --help, -h                           show this help message
  --completion                         show command completion script
  --output OUTPUT, -o OUTPUT           Path to output file. Default is stdout
  --data-yaml DATA-YAML, -y DATA-YAML  Path to yaml data file to use for templating
  --data-toml DATA-TOML, -t DATA-TOML  Path to toml data file to use for templating
  --data-env DATA-ENV, -e DATA-ENV     Path to env data file to use for templating
  --data-json DATA-JSON, -j DATA-JSON  Path to json data file to use for templating
```

## License

This project is released under the terms of the [MIT License](./LICENSE.txt).

## Contribution

### New data file format

In order to add a new data file format, you have to add a `[]byte` format decoder, a format name and a test data file as follow:

1. business code:

    Add your decoder to the list of decoders `DecodersByFormat` in `internal/decoder.go`

    ```go
    // DataDecoder fill input map data with input bytes
    type DataDecoder = func(input []byte, data Data) error

    var DecodersByFormat = map[string]DataDecoder{
      "json": jsonDecoder,
      "yaml": yamlDecoder,
      "toml": tomlDecoder,
      "env":  envDecoder,
    }
    ```

2. tests:

    1. Add test data file in `tests`

        According to `<format>` such as your new data file format name (i.e. `json`, `yaml`, `toml`, etc.)
        - file name must be `data.<format>`
        - a key `"<format>"` must equal the string `"test <format>"`
        - a key `"format"` must equal the string `"<format>"`

        (see examples in the directory `tests`)

    2. Tape `go test` in a terminal
