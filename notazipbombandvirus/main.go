package main

import (
	"fmt"
)

func main() {
	for {
		DisplayBNF()
		input := ReadInputString()
		if input == "END" {
			fmt.Println("Exiting.")
			return
		}

		node, err := ParseGraph(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			fmt.Println("Press Enter to continue...")
			fmt.Scanln()
			continue
		}

		// Generate and show leftmost derivation
		steps := LeftmostDerivation(node)
		fmt.Println("Leftmost derivation:")
		for i, s := range steps {
			fmt.Printf("%2d: %s\n", i+1, s)
		}

		fmt.Println("Derivation successful. Press Enter to display parse tree...")
		fmt.Scanln()
		// Show ASCII tree in terminal
fmt.Println("\nParse tree (terminal view):")
PrintTreeTerminal(node)

	}
}
