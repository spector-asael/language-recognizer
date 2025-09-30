
package main

import "strings"

// leftmostDerivation walks the AST in a preorder, left‑to‑right manner to
// produce the sequence of sentential forms corresponding to a leftmost
// derivation.  At each non‑terminal node, it substitutes the leftmost
// non‑terminal in the current string with the concatenation of the node's
// children.  Punctuation is preserved in the output string without
// additional spaces and adjacent <x>/<y> pairs are collapsed into a single
// token (e.g. <x><y> → <xy>).  The resulting slice contains each step
// including the final sentence.
func leftmostDerivation(root *Node) []string {
    steps := []string{}
    // Start with the root non-terminal (start symbol)
    var current string
    if root != nil {
        current = root.Value
        steps = append(steps, current)
    }
    // Helper to build children string with proper spacing and punctuation.
    buildChildrenStr := func(children []*Node) string {
        vals := make([]string, len(children))
        for i, c := range children {
            vals[i] = c.Value
        }
        var sb strings.Builder
        // Track whether to insert a space before the next token.  This
        // variable starts false for the first token.  It becomes true after
        // normal tokens and semicolons, but false after commas, since no
        // space should be inserted following a comma.
        needsSpace := false
        for i, val := range vals {
            if i == 0 {
                sb.WriteString(val)
                // Set initial spacing rule based on the token.
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
            // For subsequent tokens decide on a leading space.
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
                // Find leftmost non‑terminal (substring enclosed in <...>) and replace it.
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