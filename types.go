package main

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
