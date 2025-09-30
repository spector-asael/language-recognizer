package main

import (
	"fmt"
	"testing"
	"github.com/spector-asael/language-recognizer/parser"
	"github.com/spector-asael/language-recognizer/tree"
)

func TestParseGraphStringsVerbose(t *testing.T) {
	tests := []struct {
		input      string
		shouldPass bool
	}{
		// ✅ Valid strings
		{"HI bar A1,2 BYE", true},
		{"HI fill C3 BYE", true},
		{"HI line B2,D4 BYE", true},
		{"HI bar E5,1; fill A2 BYE", true},
		{"HI fill B1; line C3,E5 BYE", true},
		{"HI line D4,A2; bar B3,5 BYE", true},
		{"HI bar C1,4; line E2,B5; fill D3 BYE", true},
		{"HI fill A5; fill B2; fill C3 BYE", true},
		{"HI line E1,C5; bar D3,2 BYE", true},
		{"HI bar A2,4; line B3,D5 BYE", true},
		{"HI fill D1; line E2,B3; bar C4,5 BYE", true},
		{"HI line A1,B2; line C3,D4; fill E5 BYE", true},
		{"HI bar B5,1; fill D3; line A2,C4 BYE", true},
		{"HI fill E2; bar C1,5; line B3,D4 BYE", true},
		{"HI line D1,A5; bar B2,3; fill C4 BYE", true},

		// ❌ Invalid strings
		{"bar A1,2 BYE", false},
		{"HI fill F1 BYE", false},
		{"HI line B2,6 BYE", false},
		{"HI bar C3; fill D2 BYE", false},
		{"HI line A1;B2 BYE", false},
		{"HI fill A0 BYE", false},
		{"HI bar D3, BYE", false},
		{"HI line E2, D6 BYE", false},
		{"HI bar B1-2 BYE", false},
		{"HI fill BYE", false},
		{"HI line A1,B2; foo C3,D4 BYE", false},
		{"HI bar C1,2; line D3,E6 BYE", false},
		{"HI fill G1 BYE", false},
		{"HI bar A1,2 line B2,C3 BYE", false},
		{"BYE HI bar A1,2", false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			node, err := parser.LeftmostDerivation(tc.input)
			passed := (err == nil)

			if passed != tc.shouldPass {
				t.Errorf("input = %q, expected pass = %v, got err = %v", tc.input, tc.shouldPass, err)
			} else {
				t.Logf("input = %q passed as expected", tc.input)
			}

			if err == nil {
				// Print the parse tree
				t.Logf("Parse tree for input %q:", tc.input)
				parser.PrintTreeTerminal(node)

				// Print leftmost derivation
				steps := parser.PrintLeftmostDerivation(node)
				t.Logf("Leftmost derivation for input %q:", tc.input)
				for _, step := range steps {
					t.Logf("  %s", step)
				}
			} else {
				t.Logf("Error parsing input %q: %v", tc.input, err)
			}
		})
	}
}
