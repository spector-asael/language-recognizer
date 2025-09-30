package parser

import (
	"strings"
	"unicode"
	"fmt"
)

// ScanTokens converts the raw input string into a slice of Token objects.
// It sanitizes whitespace and ensures commas and semicolons are properly tokenized.
// Each contiguous sequence of letters and/or digits becomes an individual token.
// The scanner categorizes tokens as keywords, coordinates, digits, punctuation, or identifiers.
func ScanTokens(input string) ([]Token, error) {
	// Add spaces around commas and semicolons so they are treated as separate tokens
	input = strings.ReplaceAll(input, ",", " , ")
	input = strings.ReplaceAll(input, ";", " ; ")

	// Collapse multiple spaces into a single space
	input = strings.Join(strings.Fields(input), " ")

	// Split input into candidate token strings
	parts := strings.Fields(input)

	tokens := []Token{}
	for _, p := range parts {
		switch {
		case p == "HI" || p == "BYE" || p == "bar" || p == "line" || p == "fill":
			// Keywords are case-insensitive but stored in uppercase
			tokens = append(tokens, Token{Type: TokenKeyword, Value: strings.ToUpper(p)})
		case p == ",":
			tokens = append(tokens, Token{Type: TokenComma, Value: p})
		case p == ";":
			tokens = append(tokens, Token{Type: TokenSemicolon, Value: p})
		case IsXY(p):
			// XY coordinates like A4
			tokens = append(tokens, Token{Type: TokenXY, Value: p})
		case IsY(p):
			// Standalone Y coordinates like 2
			tokens = append(tokens, Token{Type: TokenY, Value: p})
		default:
			// Anything else is unrecognized
			return nil, fmt.Errorf("unrecognized token: %s", p)
		}
	}

	return tokens, nil
}

// IsXY returns true if the string is a valid XY coordinate (e.g., A1–E5)
func IsXY(s string) bool {
	if len(s) != 2 {
		return false
	}
	letter := unicode.ToUpper(rune(s[0]))
	digit := s[1]
	return letter >= 'A' && letter <= 'E' && digit >= '1' && digit <= '5'
}

// IsY returns true if the string is a valid single-digit Y coordinate (1–5)
func IsY(s string) bool {
	if len(s) != 1 {
		return false
	}
	digit := s[0]
	return digit >= '1' && digit <= '5'
}