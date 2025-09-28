package parser

import (
	"strings"
)

// AST node types
type NodeType int 

const (
	NT_GRAPH NodeType = iota
	NT_DRAW
	NT_ACTION
)

type Node struct {
	nt       NodeType
	prod     string
	term     string
	children []interface{}
	
}

// Leftmost derivation
func LeftmostDerivation(root *Node) []string {
	currentForm := []interface{}{root} // The root node
	steps := []string{"<graph>"}

	for {
		idx := firstNonterminalIndex(currentForm)
		
    	if idx == -1 {
        steps = append(steps, renderTerminals(currentForm))
        break
    }
		if idx == -1 {
			steps = append(steps, renderTerminals(currentForm))
			break
		}
		n := currentForm[idx].(*Node)
		rhs := []interface{}{}
		for _, c := range n.children {
			rhs = append(rhs, c)
		}
		newCurrent := []interface{}{}
		newCurrent = append(newCurrent, currentForm[:idx]...)
		newCurrent = append(newCurrent, rhs...)
		newCurrent = append(newCurrent, currentForm[idx+1:]...)
		currentForm = newCurrent
		steps = append(steps, renderWithNonterms(currentForm))
	}
	return steps
}

func firstNonterminalIndex(symbols []interface{}) int {
    for i, el := range symbols {
        if _, ok := el.(*Node); ok {
            return i
        }
    }
    return -1
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

