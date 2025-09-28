package parser

import (
	"strings"
)

// AST (Abstract Syntax Tree) node types
type NodeType int 

// Constants to represent different types of AST nodes
const (
	NT_GRAPH NodeType = iota // The root graph node
	NT_DRAW                  // A <draw> node
	NT_ACTION                // An action node like "bar", "line", or "fill"
)

// Node represents a node in the parse tree / AST
type Node struct {
	nt       NodeType        // The type of node
	prod     string          // The production rule string for this node
	term     string          // The terminal string (for leaf nodes)
	children []interface{}   // Children nodes or terminal strings
}

// LeftmostDerivation performs a leftmost derivation on the AST
// It expands nonterminal nodes step by step, recording each intermediate form
func LeftmostDerivation(root *Node) []string {
	// Start with the root node in the current form
	currentForm := []interface{}{root}
	
	// Steps will store the intermediate derivation strings
	steps := []string{"<graph>"} // Initial form is the root graph

	for {
		// Find the first nonterminal in the current form
		idx := firstNonterminalIndex(currentForm)
		
		// If no nonterminal remains, render the final form and break
		if idx == -1 {
			steps = append(steps, renderTerminals(currentForm))
			break
		}

		// Convert the nonterminal element to a Node
		n := currentForm[idx].(*Node)

		// Get the right-hand side (children) of this nonterminal
		rhs := []interface{}{}
		for _, c := range n.children {
			rhs = append(rhs, c)
		}

		// Replace the nonterminal in currentForm with its RHS
		newCurrent := []interface{}{}
		newCurrent = append(newCurrent, currentForm[:idx]...) // All elements before nonterminal
		newCurrent = append(newCurrent, rhs...)               // Replace with children
		newCurrent = append(newCurrent, currentForm[idx+1:]...) // All elements after nonterminal

		// Update currentForm for the next iteration
		currentForm = newCurrent

		// Render the current form, keeping nonterminals visible
		steps = append(steps, renderWithNonterms(currentForm))
	}

	return steps
}

// firstNonterminalIndex returns the index of the first nonterminal node (*Node)
// in a list of symbols. Returns -1 if there are no nonterminals.
func firstNonterminalIndex(symbols []interface{}) int {
	for i, el := range symbols {
		if _, ok := el.(*Node); ok {
			return i
		}
	}
	return -1
}

// renderWithNonterms renders a list of symbols as a string
// Nonterminal nodes are shown as their symbolic representation like "<graph>", "<draw>", "<action>"
func renderWithNonterms(cur []interface{}) string {
	parts := []string{}
	for _, el := range cur {
		switch v := el.(type) {
		case string:
			// Terminal strings are appended directly
			parts = append(parts, v)
		case *Node:
			// Nonterminals are shown as symbolic tags
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

// renderTerminals renders a list of symbols as a string
// Terminal nodes are displayed as their actual values, while nonterminals
// without a terminal are shown as "<production_rule>"
func renderTerminals(cur []interface{}) string {
	parts := []string{}
	for _, el := range cur {
		switch v := el.(type) {
		case string:
			// Terminal string
			parts = append(parts, v)
		case *Node:
			if v.term != "" {
				// Leaf node terminal value
				parts = append(parts, v.term)
			} else {
				// Nonterminal without terminal
				parts = append(parts, "<"+v.prod+">")
			}
		}
	}
	return strings.Join(parts, " ")
}
