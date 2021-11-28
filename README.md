# gotempl

Small binary used to generate files from Go Templates, environment variables and data files.

The following formats are supported:

 - JSON
 - YAML
 - TOML

## Usage

```
usage: gotempl [-h|--help] -t|--template "<value>" -d|--data "<value>"
               [-f|--format "<value>"] [-o|--output "<value>"]

               Generic templating tool

Arguments:

  -h  --help      Print help information
  -t  --template  Path to Go Template file. Example: \"TEST environment variable is {{ .Env.TEST }} and TEST data value is {{ .Data.TEST }}.\"
  -d  --data      Path to data file to use for templating
  -f  --format    Format of data file (json, yaml toml or env, defaults to json).
                  Default: json
  -o  --output    Path to output file (leave empty for stdout). Default:
```

## License

This project is released under the terms of the [MIT License](./LICENSE.txt).
