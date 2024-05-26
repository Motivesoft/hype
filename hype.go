package main

import (
	"log"
	"os"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Formatter string
	Style     string
}

func main() {
	// Check the arguments
	if len(os.Args) <= 1 {
		log.Fatal("Missing filename")
		return
	}

	// Read the file
	binary, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	// Treat the input file as a string
	code := string(binary)

	// Find out its lexer
	lexer := lexers.Analyse(code)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Load the config
	viper.SetConfigName(".hype")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.SetDefault("Formatter", "terminal16m")
	viper.SetDefault("Style", "monokai")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	
	// If there's an error or the config doesn't exist, use some defaults
	if err != nil {
		log.Fatal(err)
		return
	}
	
	// Process the configuration
	config = new(Config)
	viper.Unmarshal(config)

	// Print the file through the highlighter 
	err = quick.Highlight(os.Stdout, code, lexer.Config().Name, config.Formatter, config.Style)
	if err != nil {
		log.Fatal(err)
	}
}
