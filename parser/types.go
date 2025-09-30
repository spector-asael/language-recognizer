package parser

// Node represents a single node in the parse tree. Each node holds a value
// (either a terminal token or a non‑terminal symbol enclosed in angle brackets)
// and zero or more child nodes. The tree structure directly mirrors the
// productions applied when parsing an input string.
type Node struct {
	Value    string
	Children []*Node
}

// NewNode returns a new Node with the supplied value and no children.
func NewNode(value string) *Node {
	return &Node{Value: value, Children: make([]*Node, 0)}
}

// Append adds a new child node with the given value and returns it. This
// convenience method simplifies tree construction during parsing.
func (n *Node) Append(value string) *Node {
	child := NewNode(value)
	n.Children = append(n.Children, child)
	return child
}

// AppendNode attaches an existing node as a child. This is used when
// recursive calls in the parser return preconstructed subtrees.
func (n *Node) AppendNode(child *Node) {
	n.Children = append(n.Children, child)
}

type TokenType int

const (
    TokenKeyword TokenType = iota
    TokenXY
    TokenY
    TokenComma
    TokenSemicolon
	TokenIdentifier
)

type Token struct {
    Type  TokenType
    Value string
}