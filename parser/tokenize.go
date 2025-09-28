package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParseProgram is the entry point of the parser.
// It trims input, tokenizes it, and parses the entire program into a parse tree.
func ParseProgram(input string) (*Node, error) {
	// Step 1: Trim extra whitespace and validate input
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("empty input")
	}

	// Step 2: Tokenize input (split into meaningful tokens like "HI", "bar", coordinates, etc.)
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}

	// Step 3: Parse all the tokens found inside the <graph> 
	parsedGraph, err := parseGraphTokens(tokens)
	if err != nil {
		return nil, err
	}

	// Step 4: Build and return the root node (<graph>)
	root := &Node{
		nt:       NT_GRAPH,
		prod:     "HI <draw> BYE",
		children: []interface{}{"HI", parsedGraph, "BYE"},
	}

	return root, nil
}

// tokenize takes a raw input string and splits it into tokens.
// Tokens include keywords (HI, BYE, bar, line, fill), punctuation, or coordinates.
func tokenize(input string) ([]string, error) {
	// Step 1: Add spaces around commas and semicolons so they become separate tokens
	replacedInputString := strings.ReplaceAll(strings.ReplaceAll(input, ",", " , "), ";", " ; ")
	// Step 2: Split the string into whitespace-separated parts
	parts := strings.Fields(replacedInputString)

	// Step 3: Validate each token
	for _, p := range parts {
		// Directly allow known keywords and punctuation
		if p == "," || p == ";" || p == "HI" || p == "BYE" || p == "bar" || p == "line" || p == "fill" {
			continue
		}
		// Otherwise, check if the token matches coordinate patterns (e.g., A1, B5) or single numbers
		if !(regexp.MustCompile(`^[A-E][1-5]$`).MatchString(p) || regexp.MustCompile(`^[1-5]$`).MatchString(p)) {
			return nil, fmt.Errorf("%s contains an error - unrecognized token", p)
		}
	}

	return parts, nil
}

// parserState keeps track of our current position while parsing tokens
type parserState struct {
	tokens []string // list of tokens to parse
	pos    int      // current position in the list
}

// peek returns the current token without advancing the position
func (ps *parserState) peek() string {
	if ps.pos >= len(ps.tokens) {
		return ""
	}
	return ps.tokens[ps.pos]
}

// next returns the current token and advances the position
func (ps *parserState) next() string {
	t := ps.peek()
	if ps.pos < len(ps.tokens) {
		ps.pos++
	}
	return t
}

// parseGraphTokens parses a sequence of tokens as a <graph> node
func parseGraphTokens(tokens []string) (*Node, error) {
	// Initialize parser state
	ps := &parserState{tokens: tokens, pos: 0}

	// Parse the opening token "HI"
	firstToken := ps.next()
	if firstToken != "HI" {
		return nil, errors.New(`input must start with "HI"`)
	}

	// Parse drawing commands (<draw>)
	drawNode, err := parseDrawTokens(ps)
	if err != nil {
		return nil, err
	}

	// Parse the closing token "BYE"
	lastToken := ps.next()
	if lastToken != "BYE" {
		return nil, errors.New(`missing "BYE" or extra tokens after drawing commands`)
	}

	// Ensure no extra tokens remain after "BYE"
	if ps.pos != len(ps.tokens) {
		return nil, errors.New("extra unexpected tokens after BYE")
	}

	return drawNode, nil
}

// parseDrawTokens parses a sequence of tokens as a <draw> node
func parseDrawTokens(ps *parserState) (*Node, error) {
	// Parse a single <action>
	aNode, err := parseActionTokens(ps)
	if err != nil {
		return nil, err
	}

	// If there's a semicolon, parse another <draw> after it (recursive parsing)
	if ps.peek() == ";" {
		ps.next() // consume the semicolon
		dNode, err := parseDrawTokens(ps)
		if err != nil {
			return nil, err
		}
		return &Node{
			nt:       NT_DRAW,
			prod:     "<action> ; <draw>",
			children: []interface{}{aNode, ";", dNode},
		}, nil
	}

	// Otherwise, just return the single <action>
	return &Node{
		nt:       NT_DRAW,
		prod:     "<action>",
		children: []interface{}{aNode},
	}, nil
}

// parseActionTokens parses a sequence of tokens as a single <action> node
func parseActionTokens(ps *parserState) (*Node, error) {
	p := ps.peek()

	// "bar" action: expects format "bar <x><y> , <y>"
	if p == "bar" {
		ps.next()
		xy := ps.next()
		if xy == "" {
			return nil, errors.New("expected coordinate after 'bar'")
		}
		if ps.next() != "," {
			return nil, errors.New("expected ',' in bar action")
		}
		y2 := ps.next()
		if !isXY(xy) || !isY(y2) {
			return nil, fmt.Errorf("error: %s or %s contains an unrecognized value", xy, y2)
		}
		term := fmt.Sprintf("bar %s , %s", xy, y2)
		return &Node{
			nt:       NT_ACTION,
			prod:     "bar <x><y>,<y>",
			children: []interface{}{"bar ", xy, ",", y2},
			term:     term,
		}, nil
	} 

	// "line" action: expects format "line <x><y> , <x><y>"
	if p == "line" {
		ps.next()
		p1 := ps.next()
		if p1 == "" {
			return nil, errors.New("expected coordinate after 'line'")
		}
		if ps.next() != "," {
			return nil, errors.New("expected ',' in line action")
		}
		p2 := ps.next()
		if !isXY(p1) || !isXY(p2) {
			return nil, fmt.Errorf("error: %s or %s contains an unrecognized value", p1, p2)
		}
		term := fmt.Sprintf("line %s , %s", p1, p2)
		return &Node{
			nt:       NT_ACTION,
			prod:     "line <x><y>,<x><y>",
			children: []interface{}{"line ", p1, ",", p2},
			term:     term,
		}, nil
	} 

	// "fill" action: expects format "fill <x><y>"
	if p == "fill" {
		ps.next()
		xy := ps.next()
		if !isXY(xy) {
			return nil, fmt.Errorf("error: %s contains an unrecognized value", xy)
		}
		term := fmt.Sprintf("fill %s", xy)
		return &Node{
			nt:       NT_ACTION,
			prod:     "fill <x><y>",
			children: []interface{}{"fill ", xy},
			term:     term,
		}, nil
	}

	// If no recognized action, return an error
	return nil, fmt.Errorf("error: action '%s' not valid", p)
}

// isXY validates coordinates in the form "A1" to "E5"
func isXY(token string) bool {
	return regexp.MustCompile(`^[A-E][1-5]$`).MatchString(token)
}

// isY validates a single Y-coordinate (1-5)
func isY(token string) bool {
	return regexp.MustCompile(`^[1-5]$`).MatchString(token)
}
