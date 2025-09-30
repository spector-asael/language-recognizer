package main

import (
	"fmt"
	"strings"
	"github.com/spector-asael/language-recognizer/parser"
)

// main orchestrates the program: it repeatedly displays the grammar,
// accepts an input string, scans and parses it, performs the leftmost
// derivation and draws the parse tree.  On encountering an error it
// reports the problem and prompts the user again.  When the user types
// "END" (case insensitive) the program terminates.
func main() {
	for { // Loop back whenever the user has inputted a string
		DisplayGrammarBNF() // Displays BNF upon starting
		input := ReadInputString() // Accepts an input string
		if strings.EqualFold(input, "END") { // Ends program if prompted to
			fmt.Println("Exiting.")
			return
		}

		tokens, err := scanTokens(input) // Tokenizes the input 
		if err != nil {
			fmt.Println("Failed to scan tokens:", err)
			fmt.Println("Press Enter to continue...") // prompts the user to press a key or click to continue
			fmt.Scanln()
		}

		node, err := parser.LeftmostDerivation(tokens) // Tokenize the input string 
		if err != nil { // If there are unrecognizeable tokens
			fmt.Printf("Error: %s\n", err.Error()) // It generates an appropriate error
			fmt.Println("Press Enter to continue...") // prompts the user to press a key or click to continue
			fmt.Scanln()
			continue // then resets the for loop
		}

		steps := parser.PrintLeftmostDerivation(node) // Generate and show leftmost derivation
		fmt.Println("Leftmost derivation:")
		for i, s := range steps {
			fmt.Printf("%2d: %s\n", i+1, s)
		}

		fmt.Println("Derivation successful. Press Enter to display parse tree...")
		fmt.Scanln() 
		fmt.Println("\nParse tree (terminal view):")
		parser.PrintParseTree(node) // Generate and show the parse tree
		fmt.Println("Press Enter to continue...")
		fmt.Scanln() 
	}
}


// displayGrammar prints the BNF rules for the drawing language recogniser.
// These lines are shown each time the program prompts for a new input.  The
// grammar defines the syntax accepted by the recogniser.
func DisplayGrammarBNF() {
	fmt.Println("BNF grammar for the language recognizer:")
	fmt.Println("<graph> -> HI <draw> BYE")
	fmt.Println("<draw>  -> <action> | <action> ; <draw>")
	fmt.Println("<action> -> bar <x><y>,<y> | line <x><y>,<x><y> | fill <x><y>")
	fmt.Println("<x> -> A | B | C | D | E")
	fmt.Println("<y> -> 1 | 2 | 3 | 4 | 5")
	fmt.Println()
}
