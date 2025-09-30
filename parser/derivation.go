package parser

import "strings"

// LeftmostDerivation walks the AST in preorder to produce leftmost derivation steps
func LeftmostDerivation(root *Node) []string {
	steps := []string{}
	var current string
	if root != nil {
		current = root.Value
		steps = append(steps, current)
	}

	buildChildrenStr := func(children []*Node) string {
		vals := make([]string, len(children))
		for i, c := range children {
			vals[i] = c.Value
		}
		var sb strings.Builder
		needsSpace := false
		for i, val := range vals {
			if i == 0 {
				sb.WriteString(val)
				switch val {
				case ",":
					needsSpace = false
				case ";":
					needsSpace = true
				default:
					needsSpace = true
				}
				continue
			}
			if needsSpace && val != "," && val != ";" {
				sb.WriteString(" ")
			}
			sb.WriteString(val)
			switch val {
			case ",":
				needsSpace = false
			case ";":
				needsSpace = true
			default:
				needsSpace = true
			}
		}
		return sb.String()
	}

	var traverse func(node *Node)
	traverse = func(node *Node) {
		if node == nil {
			return
		}
		if len(node.Children) > 0 {
			childrenStr := buildChildrenStr(node.Children)
			if childrenStr != "" {
				start, end := -1, -1
				for i := 0; i < len(current); i++ {
					if current[i] == '<' {
						start = i
						for j := i + 1; j < len(current); j++ {
							if current[j] == '>' {
								end = j
								break
							}
						}
						break
					}
				}
				if start != -1 && end != -1 {
					current = current[:start] + childrenStr + current[end+1:]
					steps = append(steps, current)
				}
			}
		}
		for _, child := range node.Children {
			traverse(child)
		}
	}
	traverse(root)
	return steps
}
