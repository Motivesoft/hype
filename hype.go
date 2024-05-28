package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"

	"github.com/spf13/viper"
)

var version string

var config *Config

type Config struct {
	Lexer     string
	Formatter string
	Style     string
	Trim      bool
}

func main() {
	// Check the arguments
	if len(os.Args) <= 1 {
		log.Fatal("Missing filename")
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("Version %s\n", version)
		return
	}

	// Load the config, looking for a config file in the install location, the home dir and the current location
	viper.SetConfigName(".hype")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Dir(os.Args[0]))
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.SetDefault("Lexer", "")
	viper.SetDefault("Formatter", "terminal16m")
	viper.SetDefault("Style", "monokai")
	viper.SetDefault("Trim", false)

	// Look for config files. If there's an error other than that the config doesn't exist, report it
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore this and just use our defaults
		} else {
			// Config file was found but another error was produced
			log.Print(err)
		}
	}

	// Process the configuration
	config = new(Config)
	viper.Unmarshal(config)

	// Find out which lexer to use - firstly by checking the filename
	lexer := lexers.Match(os.Args[1])
	if lexer == nil {
		// Use a fallback lexer specified in configuration?
		if config.Lexer != "" {
			lexer = lexers.Get(config.Lexer)
		}

		// What if we still have nothing? Treat the file as plain
		if lexer == nil {
			// Use the standard fallback lexer, but output as plain text with the "noop" formatter
			lexer = lexers.Fallback
			config.Formatter = "noop"
		}
	}

	//fmt.Printf("lexer: %s\n", lexer.Config().Name)
	//fmt.Printf("formatter: %s\n", config.Formatter)
	//fmt.Printf("style: %s\n", config.Style)

	// Read the file
	binary, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	// Treat the input file as a string
	code := string(binary)

	if config.Trim {
		code = strings.Trim(code, "\f\t\r\n ")
	}

	// Print the file through the highlighter
	err = quick.Highlight(os.Stdout, code, lexer.Config().Name, config.Formatter, config.Style)
	if err != nil {
		log.Fatal(err)
	}

	if config.Trim {
		fmt.Println()
	}
}
