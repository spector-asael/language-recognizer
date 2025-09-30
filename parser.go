package main

import "fmt"

// ParserError describes an error encountered during parsing.  It embeds a
// message and optionally the offending token for richer diagnostics.
type ParserError struct {
    Msg   string
    Token *Token
}

// Error implements the error interface for ParserError.
func (p *ParserError) Error() string {
    return fmt.Sprintf("%s", p.Msg)
}

// parseGraph builds an AST for the <graph> non‑terminal given a slice of
// tokens.  The parse functions consume tokens by advancing the position
// pointer (passed by reference).  If a syntax error is detected, a
// descriptive ParserError is returned.
func parseGraph(tokens []Token) (*Node, error) {
    pos := 0
    root := NewNode("<graph>")
    if pos >= len(tokens) || tokens[pos].Type != TokenKeyword || tokens[pos].Value != "HI" {
        return nil, &ParserError{Msg: "Error: program must start with 'HI'"}
    }
    root.Append(tokens[pos].Value)
    pos++
    drawNode, err := parseDraw(tokens, &pos)
    if err != nil {
        return nil, err
    }
    root.AppendNode(drawNode)
    if pos >= len(tokens) || tokens[pos].Type != TokenKeyword || tokens[pos].Value != "BYE" {
        return nil, &ParserError{Msg: "Error: program must end with 'BYE'"}
    }
    root.Append(tokens[pos].Value)
    pos++
    if pos < len(tokens) {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: unexpected token '%s' after 'BYE'", tokens[pos].Value)}
    }
    return root, nil
}

// parseDraw recognises one or more actions separated by semicolons.  It
// implements the production <draw> → <action> | <action> ; <draw>.  A
// <draw> node is always returned, even if there is only a single action.
func parseDraw(tokens []Token, pos *int) (*Node, error) {
    drawNode := NewNode("<draw>")
    actionNode, err := parseAction(tokens, pos)
    if err != nil {
        return nil, err
    }
    drawNode.AppendNode(actionNode)
    for *pos < len(tokens) && tokens[*pos].Type == TokenSemicolon {
        drawNode.Append(tokens[*pos].Value)
        *pos = *pos + 1
        nextDraw, err := parseDraw(tokens, pos)
        if err != nil {
            return nil, err
        }
        drawNode.AppendNode(nextDraw)
        return drawNode, nil
    }
    return drawNode, nil
}

// parseAction parses one of the built‑in actions and its required
// coordinates.  It implements the production <action> → bar <x><y>,<y> | line
// <x><y>,<x><y> | fill <x><y>.  Unknown action names are detected here and
// reported via a ParserError.  The position pointer is advanced past the
// parsed tokens.
func parseAction(tokens []Token, pos *int) (*Node, error) {
    if *pos >= len(tokens) {
        return nil, &ParserError{Msg: "Error: expected an action after 'HI' or ';'"}
    }
    tok := tokens[*pos]
    if tok.Type != TokenKeyword {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: action '%s' not valid", tok.Value), Token: &tok}
    }
    actionNode := NewNode("<action>")
    actionName := tok.Value
    if actionName != "BAR" && actionName != "LINE" && actionName != "FILL" {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: action '%s' not valid", tok.Value), Token: &tok}
    }
    actionNode.Append(actionName)
    *pos = *pos + 1
    switch actionName {
    case "BAR":
        xyNode, err := parseXY(tokens, pos)
        if err != nil {
            return nil, err
        }
        actionNode.AppendNode(xyNode)
        if *pos >= len(tokens) || tokens[*pos].Type != TokenComma {
            return nil, &ParserError{Msg: "Error: expected ',' after first coordinate in 'bar'"}
        }
        actionNode.Append(tokens[*pos].Value)
        *pos = *pos + 1
        yNode, err := parseY(tokens, pos)
        if err != nil {
            return nil, err
        }
        actionNode.AppendNode(yNode)
    case "LINE":
        xyNode1, err := parseXY(tokens, pos)
        if err != nil {
            return nil, err
        }
        actionNode.AppendNode(xyNode1)
        if *pos >= len(tokens) || tokens[*pos].Type != TokenComma {
            return nil, &ParserError{Msg: "Error: expected ',' after first coordinate in 'line'"}
        }
        actionNode.Append(tokens[*pos].Value)
        *pos = *pos + 1
        xyNode2, err := parseXY(tokens, pos)
        if err != nil {
            return nil, err
        }
        actionNode.AppendNode(xyNode2)
    case "FILL":
        xyNode, err := parseXY(tokens, pos)
        if err != nil {
            return nil, err
        }
        actionNode.AppendNode(xyNode)
    }
    return actionNode, nil
}

// parseXY parses a coordinate of the form <x><y>.  The scanner records
// coordinates as TokenXY tokens.  This function verifies that the letter and
// digit both fall in the permitted ranges and constructs a subtree
// representing <xy> → <x> <y>.
func parseXY(tokens []Token, pos *int) (*Node, error) {
    if *pos >= len(tokens) {
        return nil, &ParserError{Msg: "Error: expected coordinate but reached end of input"}
    }
    tok := tokens[*pos]
    if tok.Type != TokenXY {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: expected coordinate but found '%s'", tok.Value), Token: &tok}
    }
    if len(tok.Value) != 2 {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: invalid coordinate '%s'", tok.Value), Token: &tok}
    }
    letter := tok.Value[0]
    digit := tok.Value[1]
    if letter < 'A' || letter > 'E' {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: %s contains an error – variable '%c' is not valid", tok.Value, letter), Token: &tok}
    }
    if digit < '1' || digit > '5' {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: %s contains the unrecognized value %c", tok.Value, digit), Token: &tok}
    }
    xyNode := NewNode("<xy>")
    xNode := xyNode.Append("<x>")
    xNode.Append(string(letter))
    yNode := xyNode.Append("<y>")
    yNode.Append(string(digit))
    *pos = *pos + 1
    return xyNode, nil
}

// parseY parses a standalone <y> non‑terminal which matches a single digit.
// The scanner represents such tokens with TokenY.  This function ensures
// the digit lies within the allowed range 1–5 and constructs the
// corresponding subtree.
func parseY(tokens []Token, pos *int) (*Node, error) {
    if *pos >= len(tokens) {
        return nil, &ParserError{Msg: "Error: expected y‑coordinate but reached end of input"}
    }
    tok := tokens[*pos]
    if tok.Type != TokenY {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: expected y‑coordinate but found '%s'", tok.Value), Token: &tok}
    }
    if len(tok.Value) != 1 {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: invalid y‑coordinate '%s'", tok.Value), Token: &tok}
    }
    digit := tok.Value[0]
    if digit < '1' || digit > '5' {
        return nil, &ParserError{Msg: fmt.Sprintf("Error: %s contains the unrecognized value %c", tok.Value, digit), Token: &tok}
    }
    yNode := NewNode("<y>")
    yNode.Append(string(digit))
    *pos = *pos + 1
    return yNode, nil
}