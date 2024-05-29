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

###
A configuration file called `.hype` in either of the following locations:
* Executable directory
* Home directory
* Current directory

The configuration file is YAML and takes the following form:
```yaml
Style: monokai
```

Formatter and a default Lexer (in the event one cannot be inferred from the filename) can also be specified, but the need for this is unlikely.
```yaml
Formatter: terminal16m
Lexer: markdown
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
* Multiple files on the command line, output one by one