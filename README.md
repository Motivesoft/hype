# hype
Replacement for the Windows command 'type' (Linux 'cat' command) that displays text files using syntax highlighting

## Background
The application is written in Go and uses the [Chroma](github.com/alecthomas/chroma/v2) package.

## Design
The initial design is simple.
* Accept one or more filenames on the command line, 
  * or as piped input if no filenames are specified
* Load the contents into a string
* Write the content to the console using the capabilities of the `chroma` package
  * the choice of lexer, formatter and style coming from:
    * hard-coded defaults 
    * the name of the loaded file(s)
    * command line configuration switches

## Testing
```shell
go run . [filename]
```

## Building the executable
### Version numbering
Using git tags, it is possible (on Linux builds) to inject the semantic version into a build fairly automatically:
```shell
go build -ldflags "-X 'main.version=`git describe --tags`'" .
```

The following would have been required to set up the tag. The use of `git describe` here is simply to check the outcome
```shell
git tag 1.0.0
git push --tags
git describe --tags
```

### Other platforms
Build for (eg) Windows when on Linux - which is handy for the version/tag injection 
```shell
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.version=`git describe --tags`'" .
```

Run the following to get a list of valid platforms and architectures:
```shell
go tool list dist  
```

### Simple test build
Simplistic method for local testing:
```shell
 go build -ldflags "-X main.version=0.0.0" .
```

## Usage
### Windows
```shell
hype.exe [filename]
hype.exe --version
hype.exe [filename] | more 
```

### Linux / Mac
```shell
./hype [filename]
./hype --version
```

## Configuration
Configuration is controlled by command line arguments. See
```shell
./hype.exe --help
```

Configurable items include the formatter, lexer and style. The following can be used to list the available options
```shell
./hype.exe --style-list
./hype.exe --formatter-list
./hype.exe --lexer-list
```

### Formatter
* Typically, the `formatter` will not need configuration unless the output is destined for some other purpose, e.g. to have file(s) rendered as HTML.
```shell
hype --formatter html test.cpp
```

### Lexer
The `lexer` is determined by the name of the file to be displayed. Typically this is sufficient, but can be overridden on the command line if a different one is required or the input is being piped to the application
```shell
./hype.exe --lexer cpp testfile
type testfile | ./hype.exe -lexer cpp 
``` 

### Style
The `style` has a hard-coded default, but can be overridden
```shell
./hype.exe --style monokai test.cpp
```

## Suggestions
In the future, we could consider the following:
* Find a way to process the file incrementally to avoid 'large file' issues
* Full versioninfo stuff on Windows exe, with a nice icon
* Some sort of `more` or `cat` capability
* Allow new defaults to be specified/saved
