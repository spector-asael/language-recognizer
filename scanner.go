package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/spector-asael/language-recognizer/parser"
)

func ReadInputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter input string (or END to quit): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	// Makes a standard spacing
	text = strings.Join(strings.Fields(text), " ")
	return text
}
// tokenize splits the input string into valid tokens.
func scanTokens(input string) ([]string, error) {
	// Add spaces around commas and semicolons
	sanitizedInput := strings.ReplaceAll(strings.ReplaceAll(input, ",", " , "), ";", " ; ")
	parts := strings.Fields(sanitizedInput)

	// Validate tokens
	for _, tok := range parts {
		if tok == "," || tok == ";" || tok == "HI" || tok == "BYE" || tok == "bar" || tok == "line" || tok == "fill" {
			continue
		}
		if !(parser.IsXY(tok) || parser.IsY(tok)) {
			return nil, fmt.Errorf("unrecognized token: %s", tok)
		}
	}

	return parts, nil
}