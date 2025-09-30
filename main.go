package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// displayGrammar prints the BNF rules for the drawing language recogniser.
// These lines are shown each time the program prompts for a new input.  The
// grammar defines the syntax accepted by the recogniser.
func displayGrammar() {
	fmt.Println("Grammar for drawing program:")
	fmt.Println("<graph>    → HI <draw> BYE")
	fmt.Println("<draw>     → <action> | <action> ; <draw>")
	fmt.Println("<action>   → bar <x><y>,<y> | line <x><y>,<x><y> | fill <x><y>")
	fmt.Println("<x>        → A | B | C | D | E")
	fmt.Println("<y>        → 1 | 2 | 3 | 4 | 5")
	fmt.Println()
}

// main orchestrates the program: it repeatedly displays the grammar,
// accepts an input string, scans and parses it, performs the leftmost
// derivation and draws the parse tree.  On encountering an error it
// reports the problem and prompts the user again.  When the user types
// "END" (case insensitive) the program terminates.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		displayGrammar()
		fmt.Print("Enter drawing string (or END to quit): ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if strings.EqualFold(line, "END") {
			break
		}
		tokens, err := scanTokens(line)
		if err != nil {
			// On scanning error, report the problem and wait for the user
			// to press Enter before returning to the top of the loop.  This
			// behaviour follows the assignment guideline to prompt the user
			// to continue after an unsuccessful derivation attempt.
			fmt.Println(err)
			fmt.Print("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')
			continue
		}
		ast, perr := parseGraph(tokens)
		if perr != nil {
			// On parsing error, report the problem and wait for user
			fmt.Println(perr.Error())
			fmt.Print("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')
			continue
		}
		// Successful derivation: print each sentential form and prompt
		// the user before drawing the parse tree, as required by the
		// assignment.  After displaying the tree, prompt again before
		// looping back for another input.
		fmt.Println("Leftmost derivation:")
		steps := leftmostDerivation(ast)
		for i, s := range steps {
			fmt.Printf("%02d → %s\n", i+1, s)
		}
		fmt.Println()
		// Indicate success and wait for the user before drawing the tree.
		fmt.Print("SUCCESS. \n")
		fmt.Print("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println()
		fmt.Println("Parse tree:")

		printArrayTree(ast)

		fmt.Println()
		// After displaying the tree, wait again before prompting for a new string.
		fmt.Print("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
}
