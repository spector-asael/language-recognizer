package parser

import (
	"errors"
	"fmt"
	"regexp"
)
// parserState keeps track of the current token position while parsing.
type parserState struct {
	tokens []string
	pos    int
}

func (ps *parserState) peek() string {
	if ps.pos >= len(ps.tokens) {
		return ""
	}
	return ps.tokens[ps.pos]
}

func (ps *parserState) next() string {
	tok := ps.peek()
	if ps.pos < len(ps.tokens) {
		ps.pos++
	}
	return tok
}

// parseGraphTokens parses a sequence of tokens as a <graph> node.
func ParseGraphTokens(tokens []string) (*Node, error) {
	ps := &parserState{tokens: tokens, pos: 0}

	if ps.next() != "HI" {
		return nil, errors.New(`program must start with "HI"`)
	}

	if tokens[len(tokens)-1] != "BYE" {
		return nil, errors.New(`program must end with "BYE"`)
	}

	drawNode, err := parseDrawTokens(ps)
	if err != nil {
		return nil, err
	}

	if ps.next() != "BYE" {
		return nil, errors.New(`missing "BYE" or extra tokens after drawing commands`)
	}

	if ps.pos != len(ps.tokens) {
		return nil, errors.New("unexpected extra tokens after BYE")
	}

	return drawNode, nil
}

// parseDrawTokens parses a <draw> node (may be recursive for multiple actions).
func parseDrawTokens(ps *parserState) (*Node, error) {
	actionNode, err := parseActionTokens(ps)
	if err != nil {
		return nil, err
	}

	if ps.peek() == ";" {
		ps.next() // consume semicolon
		nextDrawNode, err := parseDrawTokens(ps)
		if err != nil {
			return nil, err
		}
		return &Node{
			nodeType:      NT_DRAW,
			productionRule: "<action> ; <draw>",
			children:      []interface{}{actionNode, ";", nextDrawNode},
		}, nil
	}

	return &Node{
		nodeType:      NT_DRAW,
		productionRule: "<action>",
		children:      []interface{}{actionNode},
	}, nil
}

// parseActionTokens parses a single <action> node.
func parseActionTokens(ps *parserState) (*Node, error) {
	token := ps.peek()

	switch token {
	case "BYE":
		if(token == "BYE"){
		return nil, fmt.Errorf("Your program does not contain any actions.")
	}
	case "bar":
		ps.next()
		start := ps.next()
		if start == "" {
			return nil, errors.New("expected coordinate after 'bar'")
		}
		if ps.next() != "," {
			return nil, errors.New("expected ',' in 'bar' action")
		}
		end := ps.next()
		if !isXY(start) || !isY(end) {
			return nil, fmt.Errorf("invalid coordinates: %s, %s", start, end)
		}
		term := fmt.Sprintf("bar %s , %s", start, end)
		return &Node{
			nodeType:      NT_ACTION,
			productionRule: "bar <x><y>,<y>",
			children:      []interface{}{"bar ", start, ",", end},
			terminalValue: term,
		}, nil

	case "line":
		ps.next()
		start := ps.next()
		if start == "" {
			return nil, errors.New("expected coordinate after 'line'")
		}
		if ps.next() != "," {
			return nil, errors.New("expected ',' in 'line' action")
		}
		end := ps.next()
		if !isXY(start) || !isXY(end) {
			return nil, fmt.Errorf("invalid coordinates: %s, %s", start, end)
		}
		term := fmt.Sprintf("line %s , %s", start, end)
		return &Node{
			nodeType:      NT_ACTION,
			productionRule: "line <x><y>,<x><y>",
			children:      []interface{}{"line ", start, ",", end},
			terminalValue: term,
		}, nil

	case "fill":
		ps.next()
		coord := ps.next()
		if !isXY(coord) {
			return nil, fmt.Errorf("invalid coordinate: %s", coord)
		}
		term := fmt.Sprintf("fill %s", coord)
		return &Node{
			nodeType:      NT_ACTION,
			productionRule: "fill <x><y>",
			children:      []interface{}{"fill ", coord},
			terminalValue: term,
		}, nil
	}

	return nil, fmt.Errorf("invalid action: '%s'", token)
}

// isXY validates coordinates like "A1" to "E5".
func isXY(token string) bool {
	return regexp.MustCompile(`^[A-E][1-5]$`).MatchString(token)
}

// isY validates a single Y-coordinate (1-5).
func isY(token string) bool {
	return regexp.MustCompile(`^[1-5]$`).MatchString(token)
}
