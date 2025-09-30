package parser

import (
	"fmt"
	"strings"
)

// NodePosition stores the calculated position for a node in the tree
type NodePosition struct {
	node   interface{} // can be string or *Node
	x      int
	y      int
	width  int
	center int
}

// getTreeLabel returns symbolic label for tree display
func getTreeLabel(n interface{}) string {
	switch v := n.(type) {
	case string:
		return v
	case *Node:
		switch v.nodeType {
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
		default:
			return "<unknown>"
		}
	}
	return "<unknown>"
}

// getChildren returns children of a node
func getChildren(n interface{}) []interface{} {
	switch v := n.(type) {
	case *Node:
		return v.children
	}
	return nil
}

// calculateArrayPositions performs post-order traversal to assign positions
func calculateArrayPositions(root interface{}) ([]NodePosition, int, int) {
	var positions []NodePosition
	maxDepth := 0

	var traverse func(n interface{}, depth int) (int, int)
	traverse = func(n interface{}, depth int) (int, int) {
		if n == nil {
			return 0, 0
		}

		if depth > maxDepth {
			maxDepth = depth
		}

		label := getTreeLabel(n)
		children := getChildren(n)

		if len(children) == 0 {
			nodeWidth := len(label)
			center := nodeWidth / 2
			positions = append(positions, NodePosition{node: n, y: depth, width: nodeWidth, center: center})
			return nodeWidth, center
		}

		childWidths := make([]int, len(children))
		childCenters := make([]int, len(children))
		totalWidth := 0
		gap := 4

		for i, child := range children {
			w, c := traverse(child, depth+1)
			childWidths[i] = w
			childCenters[i] = c
			if i > 0 {
				totalWidth += gap
			}
			totalWidth += w
		}

		// Center node over children
		left := childCenters[0]
		right := 0
		for i := range childWidths {
			if i > 0 {
				right += gap
			}
			if i == len(childWidths)-1 {
				right += childCenters[i]
			} else {
				right += childWidths[i]
			}
		}
		nodeCenter := (left + right) / 2
		nodeWidth := len(label)
		minWidth := nodeCenter + (nodeWidth+1)/2
		if totalWidth < minWidth {
			totalWidth = minWidth
		}
		if nodeCenter-nodeWidth/2 < 0 {
			shift := nodeWidth/2 - nodeCenter
			nodeCenter += shift
			totalWidth += shift
		}

		positions = append(positions, NodePosition{node: n, y: depth, width: totalWidth, center: nodeCenter})
		return totalWidth, nodeCenter
	}

	totalWidth, _ := traverse(root, 0)
	assignAbsolutePositions(root, 0, positions)
	return positions, totalWidth, maxDepth
}

// assignAbsolutePositions sets absolute x positions
func assignAbsolutePositions(n interface{}, xOffset int, positions []NodePosition) {
	var nodePos *NodePosition
	for i := range positions {
		if positions[i].node == n {
			nodePos = &positions[i]
			break
		}
	}
	if nodePos == nil {
		return
	}
	nodePos.x = xOffset + nodePos.center

	children := getChildren(n)
	if len(children) == 0 {
		return
	}
	childX := xOffset
	gap := 4
	for _, child := range children {
		var cp *NodePosition
		for i := range positions {
			if positions[i].node == child {
				cp = &positions[i]
				break
			}
		}
		if cp != nil {
			assignAbsolutePositions(child, childX, positions)
			childX += cp.width + gap
		}
	}
}

// PrintParseTree renders the tree in terminal
func PrintParseTree(root *Node) {
	positions, totalWidth, maxDepth := calculateArrayPositions(root)

	height := maxDepth*2 + 1
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, totalWidth+10)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for _, pos := range positions {
		row := pos.y * 2
		label := getTreeLabel(pos.node)
		startX := pos.x - len(label)/2
		for i, ch := range label {
			if startX+i >= 0 && startX+i < len(grid[row]) {
				grid[row][startX+i] = ch
			}
		}

		children := getChildren(pos.node)
		if len(children) > 0 {
			connectorRow := row + 1
			childPositions := []NodePosition{}
			for _, child := range children {
				for _, cp := range positions {
					if cp.node == child {
						childPositions = append(childPositions, cp)
						break
					}
				}
			}
			if len(childPositions) == 1 {
				grid[connectorRow][childPositions[0].x] = '│'
			} else {
				leftmost := childPositions[0].x
				rightmost := childPositions[len(childPositions)-1].x
				nodeStart := pos.x - len(label)/2
				nodeEnd := nodeStart + len(label) - 1
				for x := leftmost; x <= rightmost; x++ {
					if x < nodeStart || x > nodeEnd {
						grid[row][x] = '_'
					}
				}
				for _, cp := range childPositions {
					if cp.x < pos.x {
						grid[connectorRow][cp.x] = '/'
					} else if cp.x > pos.x {
						grid[connectorRow][cp.x] = '\\'
					} else {
						grid[connectorRow][cp.x] = '│'
					}
				}
			}
		}
	}

	for _, row := range grid {
		fmt.Println(strings.TrimRight(string(row), " "))
	}
}
