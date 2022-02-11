# gotempl

Small binary used to generate files from Go Templates, environment variables
and data files.

The following formats are supported:

- JSON
- YAML
- TOML
- ENV (environment variables in the form "key=value")

## Usage

```bash
usage: gotempl [--help] [--completion] [--template TEMPLATE] [--output OUTPUT] [--data-json DATA-JSON [DATA-JSON ...]] [--data-yaml DATA-YAML [DATA-YAML ...]] [--data-toml DATA-TOML [DATA-TOML ...]] [--data-env DATA-ENV [DATA-ENV ...]]

Generic templating tool which use both environment variables and data files as template data

optional arguments:
  --help, -h                           show this help message
  --completion                         show command completion script
  --template TEMPLATE, -t TEMPLATE     Path to Go Template file. Default is stdin.
  --output OUTPUT, -o OUTPUT           Path to output file. Default is stdout
  --data-json DATA-JSON, -j DATA-JSON  Path to JSON file
  --data-yaml DATA-YAML, -y DATA-YAML  Path to YAML file
  --data-toml DATA-TOML, -T DATA-TOML  Path to TOML file
  --data-env DATA-ENV, -e DATA-ENV     Path to ENV file
```

### Example: Rendering JSON

Let's create a file named `sample.json` containing the following:

```json
{
  "name": "John Smith"
}
```

And a file named `greeting.template` containing the following:

```tmpl
Hello {{ .Data.name }}
```

Using **gotempl**, you can then render those files to:

```bash
$ gotempl -t greeting.template --data-json sample.json
Hello John Smith
```

### Example: Rendering multiple files

Let's create a file named `greeting.yaml` containing the following:

```yaml
greeting:
  polite: Good morning
  informal: Hi
```

Then a file `sample.json` containing the following:

```json
{
  "informal": true,
  "name": "John"
}
```

And finally, a file `greeting.template` containing the following:

```tmpl
{{- if .Data.informal -}}
{{ .Data.greeting.informal }} {{ .Data.name }}
{{- else -}}
{{ .Data.greeting.polite }} {{ .Data.name }}
{{- end -}}
```

Using **gotempl**, you can then render those files to:

```bash
$ gotempl -t greeting.template --data-yaml greeting.yaml --data-json sample.json
Hi John
```

### Example: Reading template from stdin

Using the previous `sample.json` file and **gotempl**, you can render it to:

```bash
$ export GREETING="Hello"
$ cat <<EOF | gotempl --data-json sample.json
{{ .Env.GREETING }} {{ .Data.name }}
EOF
Hello John
```

### Example: Using Sprig

**gotempl** supports [Sprig functions](http://masterminds.github.io/sprig/):

```bash
$ export STR="hello"
$ cat <<EOF | gotempl
{{ .Env.STR | upper | repeat 5 }}
EOF
HELLOHELLOHELLOHELLOHELLO
```

### Example: Using local file functions

**gotempl** supports read file functions such as

| names | description |
|-|-|
| isDir, osIsDir | test if input path is a directory |
| readDir, osReadDir | returns input path files |
| readFile, osReadFile | returns input path file content |
| walkDir, osWalkDir | returns recursively input path files and directory content |
| fileExists, osFileExists | true if input file path exists |

```bash
$ cat <<EOF | gotempl
{{ isDir "." }}
EOF
true
```

## Contributing

To add a new supported format, you'll need to implement the following interface:

```go
type Decoder {
  Format()        string        // return the format name, used for the --data-*** option
  Shortcut()      string        // return the shortcut option to use
  Decode([]byte)  (Data, error) // unmarhsal the input data
}
```

The add your implementation to the decoder list in `internal/decoder/main.go`:

```go
var decoders = []Decoder{
  //...
  MyDecoder{},
}
```

## License

This project is released under the terms of the [MIT License](./LICENSE.txt).
