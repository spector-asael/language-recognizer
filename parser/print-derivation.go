package parser

import (
	"strings"
)

// AST (Abstract Syntax Tree) node types
type NodeType int 

const (
	NT_GRAPH NodeType = iota // The root graph node
	NT_DRAW                  // A <draw> node
	NT_ACTION                // An action node like "bar", "line", or "fill"
)

// Node represents a node in the parse tree / AST
type Node struct {
	nodeType       NodeType        // The type of node
	productionRule string          // The production rule string for this node
	terminalValue  string          // The terminal string (for leaf nodes)
	children       []interface{}   // Children nodes or terminal strings
}

// PrintLeftmostDerivation performs a leftmost derivation on the AST
// It expands nonterminal nodes step by step, recording each form
func PrintLeftmostDerivation(rootNode *Node) []string {
	// Start with the root node in the current form
	currentFormSymbols := []interface{}{rootNode}
	
	// Steps will store the intermediate derivation strings
	derivationSteps := []string{"<graph>"} // Initial form is the root graph

	for {
		// Find the first nonterminal in the current form
		firstNonterminalPosition := findFirstNonterminalIndex(currentFormSymbols)
		
		// If no nonterminal remains, render the final form and break
		if firstNonterminalPosition == -1 {
			derivationSteps = append(derivationSteps, renderTerminals(currentFormSymbols))
			break
		}

		// Convert the nonterminal element to a Node
		nonterminalNode := currentFormSymbols[firstNonterminalPosition].(*Node)

		// Get the right-hand side (children) of this nonterminal
		childrenNodes := []interface{}{}
		for _, child := range nonterminalNode.children {
			childrenNodes = append(childrenNodes, child)
		}

		// Replace the nonterminal in currentFormSymbols with its children
		newCurrentForm := []interface{}{}
		newCurrentForm = append(newCurrentForm, currentFormSymbols[:firstNonterminalPosition]...) // Elements before nonterminal
		newCurrentForm = append(newCurrentForm, childrenNodes...)                                // Replace with children
		newCurrentForm = append(newCurrentForm, currentFormSymbols[firstNonterminalPosition+1:]...) // Elements after nonterminal

		// Update currentFormSymbols for the next iteration
		currentFormSymbols = newCurrentForm

		// Render the current form, keeping nonterminals visible
		derivationSteps = append(derivationSteps, renderWithNonterminals(currentFormSymbols))
	}

	return derivationSteps
}

// findFirstNonterminalIndex returns the index of the first nonterminal node (*Node)
// in a list of symbols. Returns -1 if there are no nonterminals.
func findFirstNonterminalIndex(symbolList []interface{}) int {
	for position, symbol := range symbolList {
		if _, isNode := symbol.(*Node); isNode {
			return position
		}
	}
	return -1
}

// renderWithNonterminals renders a list of symbols as a string
// Nonterminal nodes are shown as their symbolic representation like "<graph>", "<draw>", "<action>"
func renderWithNonterminals(symbolList []interface{}) string {
	renderedParts := []string{}
	for _, symbol := range symbolList {
		switch typedSymbol := symbol.(type) {
		case string:
			// Terminal strings are appended directly
			renderedParts = append(renderedParts, typedSymbol)
		case *Node:
			// Nonterminals are shown as symbolic tags
			switch typedSymbol.nodeType {
			case NT_GRAPH:
				renderedParts = append(renderedParts, "<graph>")
			case NT_DRAW:
				renderedParts = append(renderedParts, "<draw>")
			case NT_ACTION:
				renderedParts = append(renderedParts, "<action>")
			}
		}
	}
	return strings.Join(renderedParts, " ")
}

// renderTerminals renders a list of symbols as a string
// Terminal nodes are displayed as their actual values, while nonterminals
// without a terminal are shown as "<production_rule>"
func renderTerminals(symbolList []interface{}) string {
	renderedParts := []string{}
	for _, symbol := range symbolList {
		switch typedSymbol := symbol.(type) {
		case string:
			// Terminal string
			renderedParts = append(renderedParts, typedSymbol)
		case *Node:
			if typedSymbol.terminalValue != "" {
				// Leaf node terminal value
				renderedParts = append(renderedParts, typedSymbol.terminalValue)
			} else {
				// Nonterminal without terminal
				renderedParts = append(renderedParts, "<"+typedSymbol.productionRule+">")
			}
		}
	}
	return strings.Join(renderedParts, " ")
}
