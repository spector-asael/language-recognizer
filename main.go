package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spector-asael/language-recognizer/parser"
)

// ReadInputString reads a line of input from the user and trims whitespace.
func ReadInputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter input string (or END to quit): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	// Normalize spacing
	text = strings.Join(strings.Fields(text), " ")
	return text
}

// DisplayGrammarBNF prints the BNF rules for the drawing language recognizer.
// These lines are shown each time the program prompts for a new input.
func DisplayGrammarBNF() {
	fmt.Println("BNF grammar for the language recognizer:")
	fmt.Println("<graph> -> HI <draw> BYE")
	fmt.Println("<draw>  -> <action> | <action> ; <draw>")
	fmt.Println("<action> -> bar <x><y>,<y> | line <x><y>,<x><y> | fill <x><y>")
	fmt.Println("<x> -> A | B | C | D | E")
	fmt.Println("<y> -> 1 | 2 | 3 | 4 | 5")
	fmt.Println()
}

func main() {
	for {
		DisplayGrammarBNF()               // Show grammar
		input := ReadInputString()        // Read user input
		if strings.EqualFold(input, "END") {
			fmt.Println("Exiting.")
			return
		}

		// Tokenize using parser's ScanTokens
		tokens, err := parser.ScanTokens(input)
		if err != nil {
			fmt.Println("Failed to scan tokens:", err)
			fmt.Println("Press Enter to continue...")
			fmt.Scanln()
			continue
		}

		// Build parse tree
		node, err := parser.ParseGraph(tokens)
		if err != nil {
			fmt.Printf("Error parsing input: %s\n", err.Error())
			fmt.Println("Press Enter to continue...")
			fmt.Scanln()
			continue
		}

		// Compute leftmost derivation
		steps := parser.LeftmostDerivation(node)
		fmt.Println("Leftmost derivation:")
		for i, s := range steps {
			fmt.Printf("%2d: %s\n", i+1, s)
		}

		fmt.Println("Derivation successful. Press Enter to display parse tree...")
		fmt.Scanln()

		// Print parse tree
		fmt.Println("\nParse tree (terminal view):")
		parser.PrintParseTree(node)

		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
	}
}
