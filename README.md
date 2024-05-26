# hype
Replacement for the Windows command 'type' that displays text files using syntax highlighting

## Background
The application is written in Go and uses the [Chroma](github.com/alecthomas/chroma/v2) package.

## Design
The initial design is simple.
* Accept a filename on the command line
* Load it into a string
* Pass it to the Lexers module to find the appropriate one
* Use the standard formatter for coloured console output
* Hard code (initially, at least) the style
* Write the file to the console with the lexer, formatter and style

## Testing
```shell
go run . [filename]
```

## Building into an executable
```shell
go build .
```

## Usage
```shell
./hype.exe [filename]
```
