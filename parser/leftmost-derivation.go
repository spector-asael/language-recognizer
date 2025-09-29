package parser

import (
	"errors"
	"fmt"
	"strings"
)

// LeftmostDerivation is the entry point for parsing a program into an AST.
// It trims the input, tokenizes it, and constructs the parse tree.
func LeftmostDerivation(input string) (*Node, error) {
	input = strings.TrimSpace(input) // Trim empty space 
	if input == "" {
		return nil, errors.New("input is empty") // Check if there is any actual input 
	}

	tokens, err := tokenize(input) // Tokenizes the input 
	if err != nil {
		return nil, err
	}

	drawNode, err := ParseGraphTokens(tokens) // Begins parsing, returns the final <draw> token node
	if err != nil {
		return nil, err
	}

	rootNode := &Node{ // Appends the <draw> node to a <graph> node. This contains the full derivation.
		nodeType:      NT_GRAPH,
		productionRule: "HI <draw> BYE",
		children:      []interface{}{"HI", drawNode, "BYE"},
	}

	return rootNode, nil
}

// tokenize splits the input string into valid tokens.
func tokenize(input string) ([]string, error) {
	// Add spaces around commas and semicolons
	sanitizedInput := strings.ReplaceAll(strings.ReplaceAll(input, ",", " , "), ";", " ; ")
	parts := strings.Fields(sanitizedInput)

	// Validate tokens
	for _, tok := range parts {
		if tok == "," || tok == ";" || tok == "HI" || tok == "BYE" || tok == "bar" || tok == "line" || tok == "fill" {
			continue
		}
		if !(isXY(tok) || isY(tok)) {
			return nil, fmt.Errorf("unrecognized token: %s", tok)
		}
	}

	return parts, nil
}
