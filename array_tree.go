package main

import (
	"fmt"
	"strings"
)

// NodePosition stores the calculated position for a node in the tree
type NodePosition struct {
	node   *Node
	x      int // horizontal position (column)
	y      int // vertical position (row/level)
	width  int // width of the node's subtree
	center int // center position of this node relative to its subtree start
}

// calculateArrayPositions performs a post-order traversal to calculate
// exact positions for all nodes in the tree. Returns a slice of positions
// and the total width needed.
func calculateArrayPositions(root *Node) ([]NodePosition, int, int) {
	positions := []NodePosition{}
	maxDepth := 0

	var traverse func(n *Node, depth int) (int, int)
	traverse = func(n *Node, depth int) (int, int) {
		if n == nil {
			return 0, 0
		}

		if depth > maxDepth {
			maxDepth = depth
		}

		// Base case: leaf node
		if len(n.Children) == 0 {
			nodeWidth := len(n.Value)
			center := nodeWidth / 2
			positions = append(positions, NodePosition{
				node:   n,
				x:      0, // Will be adjusted later
				y:      depth,
				width:  nodeWidth,
				center: center,
			})
			return nodeWidth, center
		}

		// Recursive case: process all children first
		childWidths := make([]int, len(n.Children))
		childCenters := make([]int, len(n.Children))
		totalWidth := 0
		gap := 4 // spacing between children

		for i, child := range n.Children {
			w, c := traverse(child, depth+1)
			childWidths[i] = w
			childCenters[i] = c
			if i > 0 {
				totalWidth += gap
			}
			totalWidth += w
		}

		// Calculate the center of this node based on children
		leftmostChildCenter := childCenters[0]
		rightmostChildCenter := totalWidth - childWidths[len(childWidths)-1] + childCenters[len(childCenters)-1]

		// Adjust for gaps
		rightmostChildCenter = 0
		for i := 0; i < len(childWidths); i++ {
			if i > 0 {
				rightmostChildCenter += gap
			}
			if i == len(childWidths)-1 {
				rightmostChildCenter += childCenters[i]
			} else {
				rightmostChildCenter += childWidths[i]
			}
		}

		nodeCenter := (leftmostChildCenter + rightmostChildCenter) / 2
		nodeWidth := len(n.Value)

		// Ensure the node label fits
		minWidth := nodeCenter + (nodeWidth+1)/2
		if totalWidth < minWidth {
			totalWidth = minWidth
		}

		// Also ensure left side fits
		leftNeed := nodeCenter - nodeWidth/2
		if leftNeed < 0 {
			// Shift everything right
			shift := -leftNeed
			nodeCenter += shift
			totalWidth += shift
		}

		positions = append(positions, NodePosition{
			node:   n,
			x:      0, // Will be adjusted later
			y:      depth,
			width:  totalWidth,
			center: nodeCenter,
		})

		return totalWidth, nodeCenter
	}

	totalWidth, _ := traverse(root, 0)

	// Second pass: assign absolute x positions
	assignAbsolutePositions(root, 0, positions)

	return positions, totalWidth, maxDepth
}

// assignAbsolutePositions performs a second traversal to set absolute x coordinates
func assignAbsolutePositions(n *Node, xOffset int, positions []NodePosition) {
	// Find this node's position entry
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

	// Position children
	if len(n.Children) > 0 {
		gap := 4
		childX := xOffset

		for i, child := range n.Children {
			// Find child's position
			var childPos *NodePosition
			for j := range positions {
				if positions[j].node == child {
					childPos = &positions[j]
					break
				}
			}

			if childPos != nil {
				assignAbsolutePositions(child, childX, positions)
				childX += childPos.width
				if i < len(n.Children)-1 {
					childX += gap
				}
			}
		}
	}
}

// printArrayTree renders the tree using a 2D character array for perfect alignment
func printArrayTree(root *Node) {
	positions, totalWidth, maxDepth := calculateArrayPositions(root)

	// Create 2D grid: each level needs 2 rows (node + connector)
	height := maxDepth*2 + 1
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, totalWidth+10) // Extra padding
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Place all nodes in the grid
	for _, pos := range positions {
		row := pos.y * 2
		nodeValue := pos.node.Value
		startX := pos.x - len(nodeValue)/2

		// Place node value
		for i, ch := range nodeValue {
			if startX+i >= 0 && startX+i < len(grid[row]) {
				grid[row][startX+i] = ch
			}
		}

		// Draw connectors to children
		if len(pos.node.Children) > 0 {
			connectorRow := row + 1

			// Find children positions
			childPositions := []NodePosition{}
			for _, child := range pos.node.Children {
				for _, cp := range positions {
					if cp.node == child {
						childPositions = append(childPositions, cp)
						break
					}
				}
			}

			if len(childPositions) > 0 {
				// Special case: single child - just draw vertical line
				if len(childPositions) == 1 {
					childX := childPositions[0].x
					if childX >= 0 && childX < len(grid[connectorRow]) {
						grid[connectorRow][childX] = '│'
					}
				} else {
					// Multiple children: draw underscores and connectors
					leftmost := childPositions[0].x
					rightmost := childPositions[len(childPositions)-1].x

					// Draw underscores between leftmost and rightmost children
					nodeStart := pos.x - len(nodeValue)/2
					nodeEnd := nodeStart + len(nodeValue) - 1

					for x := leftmost; x <= rightmost; x++ {
						if x >= 0 && x < len(grid[row]) {
							// Don't overwrite the node label
							if x < nodeStart || x > nodeEnd {
								grid[row][x] = '_'
							}
						}
					}

					// Draw connectors to each child
					for _, cp := range childPositions {
						if cp.x >= 0 && cp.x < len(grid[connectorRow]) {
							if cp.x == pos.x {
								grid[connectorRow][cp.x] = '│'
							} else if cp.x < pos.x {
								grid[connectorRow][cp.x] = '/'
							} else {
								grid[connectorRow][cp.x] = '\\'
							}
						}
					}
				}
			}
		}
	}

	// Print the grid
	for _, row := range grid {
		line := string(row)
		// Trim trailing spaces
		line = strings.TrimRight(line, " ")
		fmt.Println(line)
	}
}
