package main

import (
    "fmt"
    "strings"
    "unicode"
)

// TokenType enumerates the different kinds of lexical items recognised by the
// scanner.  Distinguishing these categories up front simplifies the parser.
type TokenType int

const (
    // Keyword tokens correspond to reserved words in the grammar (HI, BYE,
    // bar, line, fill).  The token's Value field stores the exact word.
    TokenKeyword TokenType = iota
    // TokenXY represents a two‑character coordinate consisting of a letter
    // followed by a digit (e.g. D2).  The Value field holds the original
    // string; the parser will further decompose it into <x> and <y> parts.
    TokenXY
    // TokenY represents a single digit used on its own (e.g. 5) in the
    // grammar.  This only occurs after comma in the bar production.
    TokenY
    // Punctuation tokens mark separators between actions and coordinates.
    // Semicolons separate actions within the <draw> non‑terminal.
    TokenSemicolon
    // Comma separates coordinate parts in certain productions.
    TokenComma
    // Identifier covers any sequence of letters that is not a recognised
    // keyword.  The parser will validate these further and report an
    // appropriate error if they appear in positions where only keywords are
    // permitted.
    TokenIdentifier
)

// Token holds a lexeme from the input along with its classified type.  The
// parser examines the Type field to decide which grammar production to
// apply.  The Value field retains the exact lexeme for error reporting.
type Token struct {
    Type  TokenType
    Value string
}

// scanTokens converts the raw input string into a slice of Token objects.  It
// removes whitespace and breaks the string at punctuation marks.  Each
// contiguous sequence of letters and/or digits becomes an individual token.
// The scanner categorises tokens as keywords, coordinates, digits,
// punctuation or generic identifiers.  If the scanner encounters any
// characters that cannot be part of a valid token it returns an error.
func scanTokens(input string) ([]Token, error) {
    tokens := []Token{}
    var buf strings.Builder
    flush := func() error {
        if buf.Len() == 0 {
            return nil
        }
        lexeme := buf.String()
        buf.Reset()
        upperLex := strings.ToUpper(lexeme)
        switch upperLex {
        case "HI", "BYE", "BAR", "LINE", "FILL":
            tokens = append(tokens, Token{Type: TokenKeyword, Value: upperLex})
            return nil
        }
        if len(lexeme) == 2 {
            r0 := rune(lexeme[0])
            r1 := rune(lexeme[1])
            if unicode.IsLetter(r0) && unicode.IsDigit(r1) {
                letter := unicode.ToUpper(r0)
                digit := r1
                if letter >= 'A' && letter <= 'E' && digit >= '1' && digit <= '5' {
                    tokens = append(tokens, Token{Type: TokenXY, Value: string([]rune{letter, digit})})
                    return nil
                }
                tokens = append(tokens, Token{Type: TokenXY, Value: string([]rune{letter, digit})})
                return nil
            }
        }
        if len(lexeme) == 1 && unicode.IsDigit(rune(lexeme[0])) {
            digit := lexeme[0]
            tokens = append(tokens, Token{Type: TokenY, Value: string(digit)})
            return nil
        }
        tokens = append(tokens, Token{Type: TokenIdentifier, Value: lexeme})
        return nil
    }
    for _, ch := range input {
        switch {
        case unicode.IsSpace(ch):
            if err := flush(); err != nil {
                return nil, err
            }
        case ch == ';':
            if err := flush(); err != nil {
                return nil, err
            }
            tokens = append(tokens, Token{Type: TokenSemicolon, Value: ";"})
        case ch == ',':
            if err := flush(); err != nil {
                return nil, err
            }
            tokens = append(tokens, Token{Type: TokenComma, Value: ","})
        default:
            if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
                buf.WriteRune(ch)
            } else {
                return nil, fmt.Errorf("invalid character '%c' in input", ch)
            }
        }
    }
    if err := flush(); err != nil {
        return nil, err
    }
    return tokens, nil
}