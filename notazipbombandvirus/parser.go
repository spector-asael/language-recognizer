package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)


// Tokenization
func tokenize(input string) ([]string, error) {
	// Add spaces around commas and semicolons
	repl := strings.ReplaceAll(strings.ReplaceAll(input, ",", " , "), ";", " ; ")
	parts := strings.Fields(repl)

	for _, p := range parts {
		if p == "," || p == ";" || p == "HI" || p == "BYE" || p == "bar" || p == "line" || p == "fill" {
			continue
		}
		if !(regexp.MustCompile(`^[A-E][1-5]$`).MatchString(p) || regexp.MustCompile(`^[1-5]$`).MatchString(p)) {
			return nil, fmt.Errorf("%s contains an error - unrecognized token", p)
		}
	}
	return parts, nil
}

// Parser state
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
	t := ps.peek()
	if ps.pos < len(ps.tokens) {
		ps.pos++
	}
	return t
}

// Parser
func ParseGraph(input string) (*Node, error) {
	if strings.TrimSpace(input) == "" {
		return nil, errors.New("empty input")
	}
	parts, tokErr := tokenize(input)
	if tokErr != nil {
		return nil, tokErr
	}
	ps := &parserState{tokens: parts, pos: 0}

	if ps.next() != "HI" {
		return nil, errors.New("Input must start with HI")
	}

	drawNode, err := parseDraw(ps)
	if err != nil {
		return nil, err
	}

	if ps.next() != "BYE" {
		return nil, errors.New("Missing BYE or extra tokens after drawing commands")
	}
	if ps.pos != len(ps.tokens) {
		return nil, errors.New("Extra unexpected tokens after BYE")
	}

	root := &Node{nt: NT_GRAPH, prod: "HI <draw> BYE", children: []interface{}{"HI", drawNode, "BYE"}}
	return root, nil
}

func parseDraw(ps *parserState) (*Node, error) {
	aNode, err := parseAction(ps)
	if err != nil {
		return nil, err
	}
	if ps.peek() == ";" {
		ps.next()
		dNode, err := parseDraw(ps)
		if err != nil {
			return nil, err
		}
		return &Node{nt: NT_DRAW, prod: "<action> ; <draw>", children: []interface{}{aNode, ";", dNode}}, nil
	}
	return &Node{nt: NT_DRAW, prod: "<action>", children: []interface{}{aNode}}, nil
}

func parseAction(ps *parserState) (*Node, error) {
	p := ps.peek()
	if p == "bar" {
		ps.next()
		xy := ps.next()
		if xy == "" {
			return nil, errors.New("Expected coordinate after 'bar'")
		}
		if ps.next() != "," {
			return nil, errors.New("Expected ',' in bar action")
		}
		y2 := ps.next()
		if !isXY(xy) || !isY(y2) {
			return nil, fmt.Errorf("Error: %s or %s contains an unrecognized value", xy, y2)
		}
		term := fmt.Sprintf("bar %s , %s", xy, y2)
		return &Node{nt: NT_ACTION, prod: "bar <x><y>,<y>", children: []interface{}{"bar ", xy, ",", y2}, term: term}, nil
	} else if p == "line" {
		ps.next()
		p1 := ps.next()
		if p1 == "" {
			return nil, errors.New("Expected coordinate after 'line'")
		}
		if ps.next() != "," {
			return nil, errors.New("Expected ',' in line action")
		}
		p2 := ps.next()
		if !isXY(p1) || !isXY(p2) {
			return nil, fmt.Errorf("Error: %s or %s contains an unrecognized value", p1, p2)
		}
		term := fmt.Sprintf("line %s , %s", p1, p2)
		return &Node{nt: NT_ACTION, prod: "line <x><y>,<x><y>", children: []interface{}{"line ", p1, ",", p2}, term: term}, nil
	} else if p == "fill" {
		ps.next()
		xy := ps.next()
		if !isXY(xy) {
			return nil, fmt.Errorf("Error: %s contains the unrecognized value", xy)
		}
		term := fmt.Sprintf("fill %s", xy)
		return &Node{nt: NT_ACTION, prod: "fill <x><y>", children: []interface{}{"fill ", xy}, term: term}, nil
	}
	return nil, fmt.Errorf("Error: action '%s' not valid", p)
}

func isXY(token string) bool {
	return regexp.MustCompile(`^[A-E][1-5]$`).MatchString(token)
}

func isY(token string) bool {
	return regexp.MustCompile(`^[1-5]$`).MatchString(token)
}

// Leftmost derivation
func LeftmostDerivation(root *Node) []string {
	cur := []interface{}{root}
	steps := []string{"<graph>"}

	for {
		idx := -1
		for i, el := range cur {
			if _, ok := el.(*Node); ok {
				idx = i
				break
			}
		}
		if idx == -1 {
			steps = append(steps, renderTerminals(cur))
			break
		}
		n := cur[idx].(*Node)
		rhs := []interface{}{}
		for _, c := range n.children {
			rhs = append(rhs, c)
		}
		newCur := []interface{}{}
		newCur = append(newCur, cur[:idx]...)
		newCur = append(newCur, rhs...)
		newCur = append(newCur, cur[idx+1:]...)
		cur = newCur
		steps = append(steps, renderWithNonterms(cur))
	}
	return steps
}

func renderWithNonterms(cur []interface{}) string {
	parts := []string{}
	for _, el := range cur {
		switch v := el.(type) {
		case string:
			parts = append(parts, v)
		case *Node:
			switch v.nt {
			case NT_GRAPH:
				parts = append(parts, "<graph>")
			case NT_DRAW:
				parts = append(parts, "<draw>")
			case NT_ACTION:
				parts = append(parts, "<action>")
			}
		}
	}
	return strings.Join(parts, " ")
}

func renderTerminals(cur []interface{}) string {
	parts := []string{}
	for _, el := range cur {
		switch v := el.(type) {
		case string:
			parts = append(parts, v)
		case *Node:
			if v.term != "" {
				parts = append(parts, v.term)
			} else {
				parts = append(parts, "<"+v.prod+">")
			}
		}
	}
	return strings.Join(parts, " ")
}

// Print parse tree
func PrintParseTree(n *Node, indent string) {
	if n == nil {
		return
	}
	fmt.Printf("%s%s\n", indent, nodeLabel(n))
	newIndent := indent + "  "
	for _, c := range n.children {
		switch v := c.(type) {
		case string:
			fmt.Printf("%s- %s\n", newIndent, v)
		case *Node:
			PrintParseTree(v, newIndent)
		}
	}
}

func nodeLabel(n *Node) string {
	switch n.nt {
	case NT_GRAPH:
		return "<graph>"
	case NT_DRAW:
		return "<draw>"
	case NT_ACTION:
		return "<action>"
	}
	return "<node>"
}
