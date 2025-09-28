package main

import (
	"fmt"
	"github.com/spector-asael/language-recognizer/parser"
)

func main() {
	for { // Loop back whenever the user has inputted a string
		DisplayBNF() // Displays BNF upon starting
		input := ReadInputString() // Accepts an input string
		if input == "END" { // Ends program if prompted to
			fmt.Println("Exiting.")
			return
		}

		node, err := parser.ParseProgram(input) // Tokenize the input string 
		if err != nil { // If there are unrecognizeable tokens
			fmt.Printf("Error: %s\n", err.Error()) // It generates an appropriate error
			fmt.Println("Press Enter to continue...") // prompts the user to press a key or click to continue
			fmt.Scanln()
			continue // then resets the for loop
		}

		steps := parser.LeftmostDerivation(node) // Generate and show leftmost derivation
		fmt.Println("Leftmost derivation:")
		for i, s := range steps {
			fmt.Printf("%2d: %s\n", i+1, s)
		}

		fmt.Println("Derivation successful. Press Enter to display parse tree...")
		fmt.Scanln() 
		fmt.Println("\nParse tree (terminal view):")
		parser.PrintTreeTerminal(node) // Generate and show the parse tree
	}
}
