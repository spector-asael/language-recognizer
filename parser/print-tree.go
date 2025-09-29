package parser

import (
	"fmt"
)

const (
	siblingSpacing = 2
	levelRowHeight = 3
)

func getNodeLabel(n *Node) string {
	switch n.nodeType {
	case NT_GRAPH:
		return "<graph>"
	case NT_DRAW:
		return "<draw>"
	case NT_ACTION:
		return "<action>"
	case NT_X:
		return "<x>"
	case NT_Y:
		return "<y>"
	}
	return "<unknown-node>"
}

func PrintTreeTerminal(root *Node) {
	if root == nil {
		fmt.Println("<empty tree>")
		return
	}
	w, h := measureWidth(root), measureDepth(root)*levelRowHeight+1
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

	grid := make([][]rune, h)
	for i := range grid {
		grid[i] = make([]rune, w)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	render(grid, root, 0, 0)

	for _, row := range grid {
		if s := string(row); len(s) > 0 && s != string(make([]rune, len(row))) {
			fmt.Println(s)
		}
	}
}

func measureWidth(x interface{}) int {
	switch v := x.(type) {
	case string:
		if v == "" {
			return 1
		}
		return len(v)
	case *Node:
		lbl := len(getNodeLabel(v))
		if len(v.children) == 0 {
			if lbl == 0 {
				return 1
			}
			return lbl
		}
		sum := (len(v.children) - 1) * siblingSpacing
		for _, c := range v.children {
			sum += measureWidth(c)
		}
		if sum < lbl {
			return lbl
		}
		return sum
	}
	return 0
}

func measureDepth(x interface{}) int {
	switch v := x.(type) {
	case string:
		return 1
	case *Node:
		max := 0
		for _, c := range v.children {
			if d := measureDepth(c); d > max {
				max = d
			}
		}
		return 1 + max
	}
	return 0
}

func render(grid [][]rune, x interface{}, startX, y int) int {
	w := measureWidth(x)
	if w <= 0 {
		return 0
	}

	lbl := ""
	switch v := x.(type) {
	case string:
		lbl = v
	case *Node:
		lbl = getNodeLabel(v)
	}

	labelX := startX + (w-len(lbl))/2
	for i, r := range lbl {
		if y >= 0 && y < len(grid) && labelX+i < len(grid[0]) {
			grid[y][labelX+i] = r
		}
	}

	node, ok := x.(*Node)
	if !ok || len(node.children) == 0 {
		return w
	}

	parentCenter := labelX + (len(lbl)-1)/2
	childStart := startX

	for _, c := range node.children {
		cw := measureWidth(c)
		if cw < 1 {
			cw = 1
		}
		childCenter := childStart + (cw-1)/2

		if y+1 < len(grid) {
			grid[y+1][parentCenter] = '|'
		}
		if y+2 < len(grid) {
			for pos := min(parentCenter, childCenter); pos <= max(parentCenter, childCenter); pos++ {
				grid[y+2][pos] = '-'
			}
			grid[y+2][parentCenter], grid[y+2][childCenter] = '+', '+'
		}

		render(grid, c, childStart, y+3)
		childStart += cw + siblingSpacing
	}
	return w
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
