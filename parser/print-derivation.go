package parser

import (
	"strings"
)

// AST (Abstract Syntax Tree) node types
type NodeType int 

const (
    NT_GRAPH NodeType = iota // Represents the root <graph> node
    NT_DRAW // Represents a <draw> nonterminal node
    NT_ACTION // Represents an <action> nonterminal node (bar, line, fill, etc.)
    NT_X      // Represents coordinate <x>
    NT_Y      // Represents coordinate <y>
) // NT = Nonterminal constants 

// Node represents a node in the parse tree / AST
type Node struct {
	nodeType       NodeType        // The type of node
	productionRule string          // The production rule string for this node
	terminalValue  string          // The terminal string (for leaf nodes)
	children       []interface{}   // Children nodes or terminal strings
}

// leftmostDerivation walks the AST in a preorder, left‑to‑right manner to
// produce the sequence of sentential forms corresponding to a leftmost
// derivation.  At each non‑terminal node, it substitutes the leftmost
// non‑terminal in the current string with the concatenation of the node's
// children.  Punctuation is preserved in the output string without
// additional spaces and adjacent <x>/<y> pairs are collapsed into a single
// token (e.g. <x><y> → <xy>).  The resulting slice contains each step
// including the final sentence.
func PrintLeftmostDerivation(rootNode *Node) []string {
    currentFormSymbols := []interface{}{rootNode}
    derivationSteps := []string{"<graph>"} // Initial form

    for {
        firstNonterminalPosition := findFirstNonterminalIndex(currentFormSymbols)
        if firstNonterminalPosition == -1 {
            // Only append final terminal form if it's not already in the steps
            if len(derivationSteps) == 0 || derivationSteps[len(derivationSteps)-1] != renderTerminals(currentFormSymbols) {
                derivationSteps = append(derivationSteps, renderTerminals(currentFormSymbols))
            }
            break
        }

        nonterminalNode := currentFormSymbols[firstNonterminalPosition].(*Node)

        // Replace nonterminal with its children
        newCurrentForm := append([]interface{}{}, currentFormSymbols[:firstNonterminalPosition]...)
        newCurrentForm = append(newCurrentForm, nonterminalNode.children...)
        newCurrentForm = append(newCurrentForm, currentFormSymbols[firstNonterminalPosition+1:]...)
        currentFormSymbols = newCurrentForm

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
// Nonterminal nodes are shown as their symbolic representation like "<graph>", "<draw>", "<action>", "<x>", "<y>"
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
            case NT_X:
                renderedParts = append(renderedParts, "<x>")
            case NT_Y:
                renderedParts = append(renderedParts, "<y>")
            default:
                // Fallback for any unexpected node type
                renderedParts = append(renderedParts, "<"+typedSymbol.productionRule+">")
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

func isFullyTerminal(symbols []interface{}) bool {
    for _, s := range symbols {
        if _, ok := s.(*Node); ok {
            return false
        }
    }
    return true
}