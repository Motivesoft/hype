# hype
Replacement for the Windows command 'type' (Linux 'cat' command) that displays text files using syntax highlighting

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

## Building into an executable with a version number
```shell
 go build -ldflags "-X main.version=1.0.0" .
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

## Suggestions
In the future, we could consider the following:
* Simplify this right down to reduce imports
* Allow command line setting of style etc.
* Allow new defaults to be specified/saved
* Stop using the 'quick' highlight technique
* Find a way to process the file incrementally to avoid 'large file' issues
* Full versioninfo stuff on Windows exe, with a nice icon
* Some sort of 'more' capability
* Other options from 'cat'
* Seed version info from git tag, if possible on all platforms or by building for all platforms on Linux
