package parser

import (
	"fmt"
)

// PrintTreeTerminal prints the entire AST (parse tree) in the terminal using ASCII art.
// The root node is passed in, and the tree structure is printed visually.
func PrintTreeTerminal(root *Node) {
	// Start recursion with the root node.
	// "" is the initial prefix (no indentation), true means the root is considered "last" for formatting.
	printTreeRecursive(root, "", true)
	fmt.Println(" ") // Add a blank line after printing the tree for readability.
}

// Recursive helper function to print a node and its children.
// n: the current node or string to print
// prefix: indentation string to maintain tree structure visually
// isLast: whether this node is the last child of its parent (affects branch drawing)
func printTreeRecursive(n interface{}, prefix string, isLast bool) {
	var label string

	// Determine the label to print based on the type of n
	switch v := n.(type) {
	case string: // Terminal string
		label = v
	case *Node: // Node of the AST
		label = nodeLabel(v) // Get the human-readable label for the node type
	default: // Unknown type
		return
	}

	// Decide which branch character to use
	branch := "├── " // Default branch
	if isLast {
		branch = "└── " // Last child uses this branch
	}

	// Print the current node with its prefix
	fmt.Println(prefix + branch + label)

	// Prepare the new prefix for children
	newPrefix := prefix
	if isLast {
		newPrefix += "    " // No vertical line for the last child
	} else {
		newPrefix += "│   " // Keep vertical line for other children
	}

	// If n is a Node, recursively print all its children
	if node, ok := n.(*Node); ok {
		for i, c := range node.children {
			last := i == len(node.children)-1 // Determine if this child is the last
			printTreeRecursive(c, newPrefix, last)
		}
	}
}

// Returns a human-readable label for the node type
func nodeLabel(n *Node) string {
	switch n.nt {
	case NT_GRAPH:
		return "<graph>"
	case NT_DRAW:
		return "<draw>"
	case NT_ACTION:
		return "<action>"
	}
	return "<node>" // Fallback for unknown node types
}