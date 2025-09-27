package main

import (
	"fmt"
)

// PrintTreeTerminal prints the parse tree as ASCII in the terminal
func PrintTreeTerminal(root *Node) {
	printTreeRecursive(root, "", true)
	fmt.Println(" ")
}

func printTreeRecursive(n interface{}, prefix string, isLast bool) {
	var label string
	switch v := n.(type) {
	case string:
		label = v
	case *Node:
		label = nodeLabel(v)
	default:
		return
	}

	branch := "├── "
	if isLast {
		branch = "└── "
	}
	fmt.Println(prefix + branch + label)

	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	// If it's a Node, recurse through children
	if node, ok := n.(*Node); ok {
		for i, c := range node.children {
			last := i == len(node.children)-1
			printTreeRecursive(c, newPrefix, last)
		}
	}
}
