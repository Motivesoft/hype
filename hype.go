package main

import (
	"bufio"
	"fmt"
	"io"

	flag "github.com/spf13/pflag"

	"os"
	"path/filepath"

	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var version string

func main() {
	// Named arguments
	var lexer string
	var formatter string
	var style string
	var showVersion bool
	var showHelp bool
	var showStyles bool
	var showLexers bool
	var showFormatters bool

	// Nameless argument for the filename(s)
	flag.String("", "", "Input file(s)")

	// Options
	flag.StringVarP(&lexer, "lexer", "l", "", "The specific lexer to use if the default is not appropriate")
	flag.StringVarP(&formatter, "formatter", "f", "terminal16m", "The output formatter to use")
	flag.StringVarP(&style, "style", "s", "github-dark", "The theme name of the desired the output style")

	flag.BoolVarP(&showVersion, "version", "v", false, "Show version information")
	flag.BoolVarP(&showHelp, "help", "h", false, "Show this help information")
	flag.BoolVar(&showStyles, "style-list", false, "Show list of available styles (themes)")
	flag.BoolVar(&showLexers, "lexer-list", false, "Show list of available lexers")
	flag.BoolVar(&showFormatters, "formatter-list", false, "Show list of available formatters")

	// Custom error handling
	flag.Usage = func() {
		fmt.Printf("Format files for display using syntax highlighting.\n")
		fmt.Printf("\nUsage:\n %s [options] filename [filename...]\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("Options:\n")
		fmt.Printf("  -s, --style [name]     - %s (default: %s)\n", flag.Lookup("style").Usage, flag.Lookup("style").DefValue)
		fmt.Printf("  -l, --lexer [name]     - %s (default is determined from the filename)\n", flag.Lookup("lexer").Usage)
		fmt.Printf("  -v, --version          - %s\n", flag.Lookup("version").Usage)
		fmt.Printf("  -h, --help             - %s\n", flag.Lookup("help").Usage)
		fmt.Printf("\nAdvanced options:\n")
		fmt.Printf("  -f, --formatter [name] - %s (default: %s)\n", flag.Lookup("formatter").Usage, flag.Lookup("formatter").DefValue)
		fmt.Printf("  --style-list           - %s\n", flag.Lookup("style-list").Usage)
		fmt.Printf("  --lexer-list           - %s\n", flag.Lookup("lexer-list").Usage)
		fmt.Printf("  --formatter-list       - %s\n", flag.Lookup("formatter-list").Usage)
	}

	// Parse the command line
	flag.Parse()

	// Handle info flags: -version, -help and others
	if showVersion {
		fmt.Printf("%s version %s\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showStyles {
		fmt.Printf("Available styles:\n")
		for _, v := range styles.Names() {
			fmt.Println(v)
		}
		os.Exit(0)
	}

	if showLexers {
		fmt.Printf("Available lexers and aliases:\n")
		for _, v := range lexers.GlobalLexerRegistry.Lexers {
			if v.Config().Aliases != nil {
				fmt.Println(v.Config().Name, v.Config().Aliases)
			} else {
				fmt.Println(v.Config().Name)
			}
		}
		os.Exit(0)
	}

	if showFormatters {
		fmt.Printf("Available formatters:\n")
		for _, v := range formatters.Names() {
			fmt.Println(v)
		}
		os.Exit(0)
	}

	// Remaining arguments - which will be the filename(s) to hype
	args := flag.Args()

	// Check we have a filename or more to process - otherwise read from stdin
	// to allow data to be piped into the program, e.g.
	// type somefile.txt | hype --style monokai
	if len(args) == 0 {
		args = append(args, "-")
	}

	// Process each filename
	for _, filename := range args {
		// Process the file
		err := ProcessFile(filename, lexer, formatter, style)
		if err != nil {
			fmt.Println("Error processing file:", err)
			os.Exit(3)
		}
	}
}

func ProcessFile(filename string, lexerName string, formatterName string, styleName string) error {
	var scanner *bufio.Scanner

	// See if we can use the named lexer, if not use the filename, otherwise use the fallback
	// i.e. the named lexer can override the filename
	lexer := lexers.Fallback
	if lexerName == "" && filename != "-" {
		namedLexer := lexers.Match(filename)
		if namedLexer != nil {
			lexer = namedLexer
		}
	} else {
		namedLexer := lexers.Get(lexerName)
		if namedLexer != nil {
			lexer = namedLexer
		}
	}

	// See if we can use the named formatter, if not use the fallback
	formatter := formatters.Fallback
	if formatterName != "" {
		namedFormatter := formatters.Get(formatterName)
		if namedFormatter != nil {
			formatter = namedFormatter
		}
	}

	// See if we can use the named style, if not use the fallback
	style := styles.Fallback
	if styleName != "" {
		namedStyle := styles.Get(styleName)
		if namedStyle != nil {
			style = namedStyle
		}
	}

	var contents string

	// Read the file, or from stdin, and print the output through the formatting package

	if filename == "-" {
		r := io.Reader(os.Stdin)

		scanner = bufio.NewScanner(r)
		if scanner == nil {
			return fmt.Errorf("unable to read input")
		}

		// This is the default anyway
		scanner.Split(bufio.ScanLines)

		// Read all the input into a string
		var sb strings.Builder
		for scanner.Scan() {
			sb.WriteString(scanner.Text() + "\n")
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading from input: %w", err)
		}

		contents = sb.String()

		// If not specified on the command line, infer the syntax from the content
		if lexerName == "" {
			contentBasedLexer := lexers.Analyse(contents)
			if contentBasedLexer != nil {
				lexer = contentBasedLexer
			}
		}
	} else {
		// Read the whole file
		binary, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error reading from file: %w", err)
		}

		// Treat the input file as a string
		contents = string(binary)
	}

	// Tokenise and format the contents
	iterator, err := lexer.Tokenise(nil, contents)
	if err != nil {
		return err
	}

	// Output to stdout
	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		return err
	}

	// Print a new line
	if !strings.HasSuffix(contents, "\n") {
		fmt.Println()
	}

	return nil
}
