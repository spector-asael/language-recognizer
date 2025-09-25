package main

import (
	"fmt"
)

type Node struct {
	Label    string
	Children []*Node
}

func New(label string, kids ...*Node) *Node {
	return &Node{Label: label, Children: kids}
}

type Parser struct {
	toks []Token
	pos  int
}

// entry point for parsing a graph
func (p *Parser) parseGraph() (*Node, error) {
	// <graph> -> HI <draw> BYE
	if !p.accept(TokenHI) {
		return nil, p.errExpected("HI")
	}
	drawNode, err := p.parseDraw()
	if err != nil {
		return nil, err
	}
	if !p.accept(TokenBYE) {
		return nil, p.errExpected("BYE")
	}
	return New("<graph>", New("HI"), drawNode, New("BYE")), nil
}

// <draw> -> <action> | <action> ; <draw>
func (p *Parser) parseDraw() (*Node, error) {
	node := New("<draw>")
	act, err := p.parseAction()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, act)

	// optional ; <draw>
	if p.peek().Type == TokenSemi {
		p.next() // consume ;
		nextDraw, err := p.parseDraw()
		if err != nil {
			return nil, err
		}
		semi := New(";")
		node.Children = append(node.Children, semi, nextDraw)
	}
	return node, nil
}

// <action> -> bar <x><y>,<y> | line <x><y>,<x><y> | fill <x><y>
func (p *Parser) parseAction() (*Node, error) {
	t := p.peek()
	switch t.Type {
	case TokenBAR:
		p.next()
		c1, x1, y1 := p.expectCoord()
		if c1 == nil {
			return nil, p.errExpected("<x><y>")
		}
		if !p.accept(TokenComma) {
			return nil, p.errExpected(",")
		}
		var yNode *Node
		if p.peek().Type == TokenDigit {
			d := p.next()
			yNode = New(string(d.Y))
		} else if p.peek().Type == TokenCoord {
			cc := p.next()
			yNode = New(string(cc.Y)) // take only the digit
		} else {
			return nil, p.errExpected("<y>")
		}
		// Create separate <x> and <y> nodes from the coordinate
		xNode := New("<x>", New(string(x1)))
		y1Node := New("<y>", New(string(y1)))
		return New("<action>", New("bar"), xNode, y1Node, New(","), yNode), nil
	case TokenLINE:
		p.next()
		c1, x1, y1 := p.expectCoord()
		if c1 == nil {
			return nil, p.errExpected("<x><y>")
		}
		if !p.accept(TokenComma) {
			return nil, p.errExpected(",")
		}
		c2, x2, y2 := p.expectCoord()
		if c2 == nil {
			return nil, p.errExpected("<x><y>")
		}
		// Create separate <x> and <y> nodes for both coordinates
		x1Node := New("<x>", New(string(x1)))
		y1Node := New("<y>", New(string(y1)))
		x2Node := New("<x>", New(string(x2)))
		y2Node := New("<y>", New(string(y2)))
		return New("<action>", New("line"), x1Node, y1Node, New(","), x2Node, y2Node), nil
	case TokenFILL:
		p.next()
		c1, x1, y1 := p.expectCoord()
		if c1 == nil {
			return nil, p.errExpected("<x><y>")
		}
		// Create separate <x> and <y> nodes
		xNode := New("<x>", New(string(x1)))
		yNode := New("<y>", New(string(y1)))
		return New("<action>", New("fill"), xNode, yNode), nil
	default:
		return nil, p.errExpected("bar|line|fill")
	}
}

func (p *Parser) expectCoord() (*Node, rune, rune) {
	if p.peek().Type != TokenCoord {
		return nil, 0, 0
	}
	t := p.next()
	lbl := fmt.Sprintf("%c%c", t.X, t.Y)
	return New(lbl), t.X, t.Y
}

func (p *Parser) accept(tp TokenType) bool {
	if p.peek().Type == tp {
		p.next()
		return true
	}
	return false
}
func (p *Parser) next() Token { tok := p.peek(); p.pos++; return tok }
func (p *Parser) peek() Token {
	if p.pos >= len(p.toks) {
		return Token{Type: TokenEOF, Pos: p.pos}
	}
	return p.toks[p.pos]
}

func (p *Parser) errExpected(want string) error {
	got := p.peek()
	return fmt.Errorf("syntax error: expected %s near position %d (got '%s')", want, got.Pos, got.Lexeme)
}
