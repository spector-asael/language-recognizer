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

// LeftmostDerivation is the entry point for parsing a program into an AST.
// It trims the input, tokenizes it, and constructs the parse tree.
func LeftmostDerivation(input []string) (*Node, error) {

	drawNode, err := ParseGraphTokens(input) // Begins parsing, returns the final <draw> token node
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
        // If the next token is BYE, the program does not contain any actions
        return nil, fmt.Errorf("Your program does not contain any actions.")

    case "bar":
        ps.next() // consume "bar"
        start := ps.next()
        if start == "" {
            return nil, errors.New("expected coordinate after 'bar'")
        }
        if ps.next() != "," {
            return nil, errors.New("expected ',' in 'bar' action")
        }
        end := ps.next()
        if !IsXY(start) || !IsY(end) {
            return nil, fmt.Errorf("invalid coordinates: %s, %s", start, end)
        }

        // Split start coordinate into X and Y
        xPart := string(start[0])
        yPart := string(start[1])

        // Create placeholder nodes for <x> and <y>
        xNode := &Node{nodeType: NT_X, terminalValue: xPart, children: []interface{}{xPart}}
        yNode := &Node{nodeType: NT_Y, terminalValue: yPart, children: []interface{}{yPart}}
        y2Node := &Node{nodeType: NT_Y, terminalValue: end, children: []interface{}{end}}

        // Create the action node with placeholders
        return &Node{
            nodeType:      NT_ACTION,
            productionRule: "bar <x><y> , <y>",
            children:      []interface{}{"bar ", xNode, yNode, ",", y2Node},
            terminalValue: fmt.Sprintf("bar %s , %s", start, end),
        }, nil

    case "line":
        ps.next() // consume "line"
        start := ps.next()
        if start == "" {
            return nil, errors.New("expected coordinate after 'line'")
        }
        if ps.next() != "," {
            return nil, errors.New("expected ',' in 'line' action")
        }
        end := ps.next()
        if !IsXY(start) || !IsXY(end) {
            return nil, fmt.Errorf("invalid coordinates: %s, %s", start, end)
        }

        // Split start and end coordinates into X and Y
        x1Node := &Node{nodeType: NT_X, terminalValue: string(start[0]), children: []interface{}{string(start[0])}}
        y1Node := &Node{nodeType: NT_Y, terminalValue: string(start[1]), children: []interface{}{string(start[1])}}
        x2Node := &Node{nodeType: NT_X, terminalValue: string(end[0]), children: []interface{}{string(end[0])}}
        y2Node := &Node{nodeType: NT_Y, terminalValue: string(end[1]), children: []interface{}{string(end[1])}}

        // Create the action node with placeholders
        return &Node{
            nodeType:      NT_ACTION,
            productionRule: "line <x><y> , <x><y>",
            children:      []interface{}{"line ", x1Node, y1Node, ",", x2Node, y2Node},
            terminalValue: fmt.Sprintf("line %s , %s", start, end),
        }, nil

    case "fill":
        ps.next() // consume "fill"
        coord := ps.next()
        if !IsXY(coord) {
            return nil, fmt.Errorf("invalid coordinate: %s", coord)
        }

        xNode := &Node{nodeType: NT_X, terminalValue: string(coord[0]), children: []interface{}{string(coord[0])}}
        yNode := &Node{nodeType: NT_Y, terminalValue: string(coord[1]), children: []interface{}{string(coord[1])}}

        // Create the action node with placeholders
        return &Node{
            nodeType:      NT_ACTION,
            productionRule: "fill <x><y>",
            children:      []interface{}{"fill ", xNode, yNode},
            terminalValue: fmt.Sprintf("fill %s", coord),
        }, nil
    }

    return nil, fmt.Errorf("invalid action: '%s'", token)
}

// isXY validates coordinates like "A1" to "E5".
func IsXY(token string) bool {
	return regexp.MustCompile(`^[A-E][1-5]$`).MatchString(token)
}

// isY validates a single Y-coordinate (1-5).
func IsY(token string) bool {
	return regexp.MustCompile(`^[1-5]$`).MatchString(token)
}
