package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func DisplayBNF() {
	fmt.Println("BNF grammar for the language recognizer:")
	fmt.Println("<graph> -> HI <draw> BYE")
	fmt.Println("<draw>  -> <action> | <action> ; <draw>")
	fmt.Println("<action> -> bar <x><y>,<y> | line <x><y>,<x><y> | fill <x><y>")
	fmt.Println("<x> -> A | B | C | D | E")
	fmt.Println("<y> -> 1 | 2 | 3 | 4 | 5")
	fmt.Println()
}

func ReadInputString() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter input string (or END to quit): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	// Makes a standard spacing
	text = strings.Join(strings.Fields(text), " ")
	return text
}
