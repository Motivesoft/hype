package main

import (
	"log"
	"os"
	"path/filepath"

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
	lexer := lexers.Match(os.Args[1])
	if lexer == nil {
		// Use the fallback lexer? Not sure - it is likely to make a mess of the output
		//lexer = lexers.Fallback

		// Instead, as we don't understand the filetype, output the file as plain text
		os.Stdout.WriteString(code)
		return
	}

	// Load the config, looking for a config file in the install location, the home dir and the current location
	viper.SetConfigName(".hype")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Dir(os.Args[0]))
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.SetDefault("Formatter", "terminal16m")
	viper.SetDefault("Style", "github-dark")
	err = viper.ReadInConfig()

	// If there's an error or the config doesn't exist, use some defaults
	if err != nil {
		log.Printf("%s", err)
	}

	// Process the configuration
	config = new(Config)
	viper.Unmarshal(config)

	// Include "fmt" to use these debug statements
	// fmt.Printf("lexer: %s\n", lexer.Config().Name)
	// fmt.Printf("formatter: %s\n", config.Formatter)
	// fmt.Printf("style: %s\n", config.Style)

	// Print the file through the highlighter
	err = quick.Highlight(os.Stdout, code, lexer.Config().Name, config.Formatter, config.Style)
	if err != nil {
		log.Fatal(err)
	}
}
