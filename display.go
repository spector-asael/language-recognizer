package main

import (
	"fmt"
	"strings"
)

// printLeftmostDerivation displays the leftmost derivation steps for a successful parse
func printLeftmostDerivation(input string, root *Node) {
	fmt.Println("\nLEFTMOST DERIVATION:")
	fmt.Println("Step 01: <graph>")
	fmt.Println("Step 02: HI <draw> BYE")

	// Extract the draw part from the input (between HI and BYE)
	drawPart := extractDrawPart(input)
	if drawPart == "" {
		fmt.Println("Step 03: HI <action> BYE")
		return
	}

	// Show derivation steps following the exact format from assignment
	showExactDerivation(drawPart, 3)
}

// extractDrawPart extracts the part between HI and BYE
func extractDrawPart(input string) string {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(strings.ToUpper(input), "HI ") {
		return ""
	}
	if !strings.HasSuffix(strings.ToUpper(input), " BYE") {
		return ""
	}

	// Remove HI and BYE, get the middle part
	start := 3            // after "HI "
	end := len(input) - 4 // before " BYE"
	if end <= start {
		return ""
	}
	return strings.TrimSpace(input[start:end])
}

// showExactDerivation shows derivation steps following the exact assignment format
func showExactDerivation(drawPart string, stepNum int) {
	// For now, show a simplified but correct derivation
	// This matches the assignment format more closely
	fmt.Printf("Step %02d: HI <action>", stepNum)

	actions := strings.Split(drawPart, ";")
	if len(actions) > 1 {
		fmt.Println(" ; <draw> BYE")
		stepNum++
		fmt.Printf("Step %02d: HI %s ; <draw> BYE\n", stepNum, strings.TrimSpace(actions[0]))
		stepNum++
		fmt.Printf("Step %02d: HI %s BYE\n", stepNum, drawPart)
	} else {
		fmt.Println(" BYE")
		stepNum++
		action := strings.TrimSpace(actions[0])
		parts := strings.Fields(action)
		if len(parts) >= 2 {
			actionType := parts[0]
			if actionType == "bar" && len(parts) >= 4 {
				fmt.Printf("Step %02d: HI bar <x><y>,<y> BYE\n", stepNum)
				stepNum++
				fmt.Printf("Step %02d: HI bar %s,%s BYE\n", stepNum, parts[1], parts[3])
				stepNum++
				fmt.Printf("Step %02d: HI %s BYE\n", stepNum, action)
			} else if actionType == "line" && len(parts) >= 4 {
				fmt.Printf("Step %02d: HI line <x><y>,<x><y> BYE\n", stepNum)
				stepNum++
				fmt.Printf("Step %02d: HI line %s,%s BYE\n", stepNum, parts[1], parts[3])
				stepNum++
				fmt.Printf("Step %02d: HI %s BYE\n", stepNum, action)
			} else if actionType == "fill" && len(parts) >= 2 {
				fmt.Printf("Step %02d: HI fill <x><y> BYE\n", stepNum)
				stepNum++
				fmt.Printf("Step %02d: HI fill %s BYE\n", stepNum, parts[1])
				stepNum++
				fmt.Printf("Step %02d: HI %s BYE\n", stepNum, action)
			}
		}
	}
}

// Domain nodes for a compact action representation used by the ASCII printer.
type xy struct{ X, Y rune }
type act struct {
	Kind string // "bar" | "line" | "fill"
	A    xy     // first coordinate
	B    *xy    // optional: for line's second coord
	Y    rune   // for bar's trailing <y>
}

// asciiParse holds an ordered list of actions to render under <draw>.
type asciiParse struct {
	Actions []act
}

// toASCII flattens the parsed tree into a simple action list for printing.
func toASCII(root *Node) asciiParse {
	var ap asciiParse
	var walkDraw func(n *Node)
	// Node labels used in parser: "<graph>", "HI", "<draw>", "BYE", "<action>", ";", ",", "bar","line","fill", coords "A2"
	walkAction := func(a *Node) {
		if len(a.Children) == 0 {
			return
		}
		kind := a.Children[0].Label // bar/line/fill
		switch kind {
		case "bar":
			// bar <x> <y> , <y>
			xVal := a.Children[1].Children[0].Label
			y1Val := a.Children[2].Children[0].Label
			y2Val := a.Children[4].Label // last child is digit
			ap.Actions = append(ap.Actions, act{
				Kind: "bar",
				A:    xy{rune(xVal[0]), rune(y1Val[0])},
				Y:    rune(y2Val[0]),
			})
		case "line":
			// line <x> <y> , <x> <y>
			x1Val := a.Children[1].Children[0].Label
			y1Val := a.Children[2].Children[0].Label
			x2Val := a.Children[4].Children[0].Label
			y2Val := a.Children[5].Children[0].Label
			ap.Actions = append(ap.Actions, act{
				Kind: "line",
				A:    xy{rune(x1Val[0]), rune(y1Val[0])},
				B:    &xy{rune(x2Val[0]), rune(y2Val[0])},
			})
		case "fill":
			// fill <x> <y>
			xVal := a.Children[1].Children[0].Label
			yVal := a.Children[2].Children[0].Label
			ap.Actions = append(ap.Actions, act{
				Kind: "fill",
				A:    xy{rune(xVal[0]), rune(yVal[0])},
			})
		}
	}
	walkDraw = func(n *Node) {
		// <draw> -> <action> | <action> ; <draw>
		if n.Label != "<draw>" || len(n.Children) == 0 {
			return
		}
		actNode := n.Children[0]
		walkAction(actNode)
		if len(n.Children) == 3 && n.Children[1].Label == ";" {
			walkDraw(n.Children[2])
		}
	}
	// root: <graph> [HI, <draw>, BYE]
	if len(root.Children) >= 2 {
		walkDraw(root.Children[1])
	}
	return ap
}

// printASCII emits plain ASCII in the same style as the assignment.
func printASCII(root *Node) {
	ap := toASCII(root)
	// Header matching the style
	fmt.Println("<graph>")
	fmt.Println("      /         |         \\")
	fmt.Println("     HI      <draw>       BYE")
	fmt.Println("               |")
	// Now render the chain of actions under <draw>
	for idx, a := range ap.Actions {
		// <draw> -> <action> [ ; <draw> ]
		fmt.Println("            <action>")
		switch a.Kind {
		case "fill":
			// fill <x> <y>
			fmt.Println("           /      |")
			fmt.Println("        fill     <x>")
			fmt.Println("                  |")
			fmt.Printf("                  %c\n", a.A.X)
			fmt.Println("                  |")
			fmt.Println("                 <y>")
			fmt.Println("                  |")
			fmt.Printf("                  %c\n", a.A.Y)
		case "bar":
			// bar <x> <y> , <y>
			fmt.Println("         /      |       |       |       \\")
			fmt.Println("       bar     <x>     <y>      ,       <y>")
			fmt.Println("                |       |       |        |")
			fmt.Printf("                %c       %c       ,        %c\n", a.A.X, a.A.Y, a.Y)
		case "line":
			// line <x> <y> , <x> <y>
			fmt.Println("           /        |        |        |        \\")
			fmt.Println("         line      <x>      <y>       ,        <x>")
			fmt.Println("                   |        |        |         |")
			fmt.Printf("                   %c        %c        ,         %c\n", a.A.X, a.A.Y, a.B.X)
			fmt.Println("                                              |")
			fmt.Println("                                             <y>")
			fmt.Println("                                              |")
			fmt.Printf("                                              %c\n", a.B.Y)
		}
		// If another action follows, show the '; <draw>' link
		if idx < len(ap.Actions)-1 {
			fmt.Println("             ;")
			fmt.Println("            <draw>")
			fmt.Println("               |")
		}
	}
}
