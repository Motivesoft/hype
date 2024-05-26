package main

import (
	"log"
	"os"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Missing filename")
		return
	}

	binary, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	code := string(binary)

	lexer := lexers.Analyse(code)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// lexer = chroma.Coalesce(lexer)

	err = quick.Highlight(os.Stdout, code, lexer.Config().Name, "terminal16m", "monokai")
	if err != nil {
		log.Fatal(err)
	}
}
