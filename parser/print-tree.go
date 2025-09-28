package parser

import (
	"fmt"
)

// PrintTreeTerminal prints the AST (parse tree) in the terminal using ASCII art.
// rootNode: the root of the AST
func PrintTreeTerminal(rootNode *Node) {
	// Start recursive printing with the root node
	// "" is the initial indentation prefix
	// true indicates that the root is considered the "last" child for formatting
	printNodeRecursive(rootNode, "", true)
	fmt.Println() // Add a blank line after the tree for readability
}

// printNodeRecursive recursively prints a node and its children in a tree-like structure.
// nodeOrString: the current node or terminal string to print
// indentationPrefix: current indentation to show tree structure
// isLastChild: whether this node is the last child of its parent
func printNodeRecursive(nodeOrString interface{}, indentationPrefix string, isLastChild bool) {
	var nodeLabelText string

	// Determine label based on type
	switch typedNode := nodeOrString.(type) {
	case string:
		nodeLabelText = typedNode // Terminal string
	case *Node:
		nodeLabelText = getNodeLabel(typedNode) // Node type label
	default:
		return // Unknown type, skip
	}

	// Determine branch symbol
	branchSymbol := "├── "
	if isLastChild {
		branchSymbol = "└── "
	}

	// Print the node with the current indentation and branch
	fmt.Println(indentationPrefix + branchSymbol + nodeLabelText)

	// Prepare prefix for children
	newIndentationPrefix := indentationPrefix
	if isLastChild {
		newIndentationPrefix += "    " // No vertical line for last child
	} else {
		newIndentationPrefix += "│   " // Vertical line for other children
	}

	// If this is a Node, recursively print its children
	if node, ok := nodeOrString.(*Node); ok {
		for index, child := range node.children {
			isLast := index == len(node.children)-1
			printNodeRecursive(child, newIndentationPrefix, isLast)
		}
	}
}

// getNodeLabel returns a human-readable label for a node type
func getNodeLabel(node *Node) string {
	switch node.nodeType {
	case NT_GRAPH:
		return "<graph>"
	case NT_DRAW:
		return "<draw>"
	case NT_ACTION:
		return "<action>"
	}
	return "<unknown-node>" // Fallback for unexpected node types
}
