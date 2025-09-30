package main

// import "fmt"

// Node represents a single node in the parse tree.  Each node holds a value
// (either a terminal token or a non‑terminal symbol enclosed in angle brackets)
// and zero or more child nodes.  The tree structure directly mirrors the
// productions applied when parsing an input string.
type Node struct {
	Value    string
	Children []*Node
}

// NewNode returns a new Node with the supplied value and no children.
func NewNode(value string) *Node {
	return &Node{Value: value, Children: make([]*Node, 0)}
}

// Append adds a new child node with the given value and returns it.  This
// convenience method simplifies tree construction during parsing.
func (n *Node) Append(value string) *Node {
	child := NewNode(value)
	n.Children = append(n.Children, child)
	return child
}

// AppendNode attaches an existing node as a child.  This is used when
// recursive calls in the parser return preconstructed subtrees.
func (n *Node) AppendNode(child *Node) {
	n.Children = append(n.Children, child)
}

// printTree recursively prints the parse tree in an ASCII art style.  The
// function distinguishes between printing the root node and printing other
// nodes.  At the root, only the node value is printed.  For all other
// nodes, branch characters (├── or └──) are drawn to indicate sibling
// relationships, and vertical bars (│) are used to show continuing branches.
//
// Parameters:
//
//	n      – the current node to print
//	prefix – the accumulated indentation and vertical branch characters
//	isLast – whether this node is the last child of its parent
//	isRoot – whether this node is the root of the tree
// func printTree(n *Node, prefix string, isLast bool, isRoot bool) {
// 	if n == nil {
// 		return
// 	}
// 	if isRoot {
// 		// The root prints its value without any branch characters.
// 		fmt.Println(n.Value)
// 	} else {
// 		// Print the prefix, branch connector and the node's value.
// 		fmt.Print(prefix)
// 		if isLast {
// 			fmt.Print("└── ")
// 		} else {
// 			fmt.Print("├── ")
// 		}
// 		fmt.Println(n.Value)
// 	}
// 	// Prepare the prefix for children.  Append vertical lines or spaces
// 	// depending on whether this node was the last child.
// 	newPrefix := prefix
// 	if !isRoot {
// 		if isLast {
// 			newPrefix += "    "
// 		} else {
// 			newPrefix += "│   "
// 		}
// 	}
// 	// Recursively print each child.  Mark the last child so that its branch
// 	// character becomes └──.
// 	for i, child := range n.Children {
// 		last := i == len(n.Children)-1
// 		printTree(child, newPrefix, last, false)
// 	}
// }
