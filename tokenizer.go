package main

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

type Token struct {
	Type  TokenType
	Value string
}

type Tokenizer struct {
	input  string
	cursor int
}

var KEYWORDS = map[string]TokenType{
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"while":  WHILE,
	"print":  PRINT,
	"return": RETURN,
}

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_CURLY
	RIGHT_CURLY
	LEFT_BRACKET
	RIGHT_BRACKET
	EQUAL
	GREATER
	LESS
	EQUAL_EQUAL
	GREATER_EQUAL
	LESS_EQUAL
	NUMBER
	PLUS
	MINUS
	MULT
	DIV
	STR_LIT
	ID
	LET
	IF
	ELSE
	FOR
	WHILE
	PRINT
	SEMICOLON
	RETURN
	EOF
)

func (tt TokenType) String() string {
	switch tt {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_CURLY:
		return "LEFT_CURLY"
	case RIGHT_CURLY:
		return "RIGHT_CURLY"
	case LEFT_BRACKET:
		return "LEFT_BRACKET"
	case RIGHT_BRACKET:
		return "RIGHT_BRACKET"
	case EQUAL:
		return "EQUAL"
	case GREATER:
		return "GREATER"
	case LESS:
		return "LESS"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case NUMBER:
		return "NUMBER"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULT:
		return "MULT"
	case DIV:
		return "DIV"
	case STR_LIT:
		return "STR_LIT"
	case ID:
		return "ID"
	case LET:
		return "LET"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case FOR:
		return "FOR"
	case WHILE:
		return "WHILE"
	case PRINT:
		return "PRINT"
	case SEMICOLON:
		return "SEMICOLON"
	case RETURN:
		return "RETURN"
	case EOF:
		return "EOF"
	}
	return "UNKNOWN TOKEN TYPE"
}

func (tok Token) String() string {
	return fmt.Sprintf("%12s: %s", tok.Type, tok.Value)
}

func NewLeftParen() Token {
	return Token{
		Type:  LEFT_PAREN,
		Value: "(",
	}
}

func NewRightParen() Token {
	return Token{
		Type:  RIGHT_PAREN,
		Value: ")",
	}
}

func NewLeftCurly() Token {
	return Token{
		Type:  LEFT_CURLY,
		Value: "{",
	}
}

func NewRightCurly() Token {
	return Token{
		Type:  RIGHT_CURLY,
		Value: "}",
	}
}

func NewLeftBracket() Token {
	return Token{
		Type:  LEFT_BRACKET,
		Value: "[",
	}
}

func NewRightBracket() Token {
	return Token{
		Type:  RIGHT_BRACKET,
		Value: "]",
	}
}

func NewGreaterEqual() Token {
	return Token{
		Type:  GREATER_EQUAL,
		Value: ">=",
	}
}
func NewGreater() Token {
	return Token{
		Type:  GREATER,
		Value: ">",
	}
}

func NewLessEqual() Token {
	return Token{
		Type:  LESS_EQUAL,
		Value: "<=",
	}
}

func NewLess() Token {
	return Token{
		Type:  LESS,
		Value: "<",
	}
}

func NewEqualEqual() Token {
	return Token{
		Type:  EQUAL_EQUAL,
		Value: "==",
	}
}

func NewEqual() Token {
	return Token{
		Type:  EQUAL,
		Value: "=",
	}
}

func NewPlus() Token {
	return Token{
		Type:  PLUS,
		Value: "+",
	}
}

func NewMinus() Token {
	return Token{
		Type:  MINUS,
		Value: "-",
	}
}

func NewMult() Token {
	return Token{
		Type:  MULT,
		Value: "*",
	}
}

func NewDiv() Token {
	return Token{
		Type:  DIV,
		Value: "/",
	}
}

func NewNumber(n string) Token {
	return Token{
		Type:  NUMBER,
		Value: n,
	}
}

func NewStrLit(s string) Token {
	return Token{
		Type:  STR_LIT,
		Value: s,
	}
}

func NewID(id string) Token {
	return Token{
		Type:  ID,
		Value: id,
	}
}

func NewSemiColon() Token {
	return Token{
		Type:  SEMICOLON,
		Value: ";",
	}
}

func NewEOF() Token {
	return Token{
		Type:  EOF,
		Value: "EOF",
	}
}

func NewTokenizer(input string) Tokenizer {
	return Tokenizer{
		input:  input,
		cursor: 0,
	}
}

func (t *Tokenizer) Next() (Token, error) {
	if t.isEnd() {
		return NewEOF(), nil
	}
	char := t.input[t.cursor]
	for unicode.IsSpace(rune(char)) {
		t.cursor++
		if t.isEnd() {
			return NewEOF(), nil
		}
		char = t.input[t.cursor]
	}

	switch {
	case char == '(':
		t.cursor++
		return NewLeftParen(), nil
	case char == ')':
		t.cursor++
		return NewRightParen(), nil
	case char == '{':
		t.cursor++
		return NewLeftCurly(), nil
	case char == '}':
		t.cursor++
		return NewRightCurly(), nil
	case char == '[':
		t.cursor++
		return NewLeftBracket(), nil
	case char == ']':
		t.cursor++
		return NewRightBracket(), nil
	case char == '+':
		t.cursor++
		return NewPlus(), nil
	case char == '-':
		t.cursor++
		return NewMinus(), nil
	case char == '*':
		t.cursor++
		return NewMult(), nil
	case char == '/':
		t.cursor++
		return NewDiv(), nil
	case char == '>':
		{
			next := t.Peek()
			if next == '=' {
				t.cursor += 2
				return NewGreaterEqual(), nil
			} else {
				t.cursor++
				return NewGreater(), nil
			}
		}
	case char == '<':
		{
			next := t.Peek()
			if next == '=' {
				t.cursor += 2
				return NewLessEqual(), nil
			} else {
				t.cursor++
				return NewLess(), nil
			}
		}
	case char == '=':
		{
			next := t.Peek()
			if next == '=' {
				t.cursor += 2
				return NewEqualEqual(), nil
			} else {
				t.cursor++
				return NewEqual(), nil
			}
		}
	case char == '"':
		str_lit := strings.Builder{}
		for {
			t.cursor++
			if t.isEnd()  {
				break
			}
			char = t.input[t.cursor]
			if char == '"' {
				t.cursor++
				break
			}
			str_lit.WriteByte(char)
		}
		return NewStrLit(str_lit.String()), nil
		
	case unicode.IsLetter(rune(char)), char == '_':
		{
			id := t.consumeIdentifier()
			typ, ok := KEYWORDS[id]
			if ok {
				tok := Token{
					Type:  typ,
					Value: id,
				}
				return tok, nil
			}
			return NewID(id), nil

		}
	case unicode.IsDigit(rune(char)):
		{
			n, err := t.consumeNumber()
			if err != nil {
				return Token{}, fmt.Errorf("next: %s", err)
			}
			return NewNumber(n), nil
		}
	case char == ';':
		t.cursor++
		return NewSemiColon(), nil
	default:
		return Token{}, fmt.Errorf("next: invalid token '%c'", char)
	}
	panic("unreachable")
}

func (t *Tokenizer) Peek() byte {
	if t.cursor+1 >= len(t.input) {
		return 0
	}
	return t.input[t.cursor+1]
}

func (t *Tokenizer) isEnd() bool {
	return t.cursor >= len(t.input)
}

func (t *Tokenizer) validIdentifierChars() string {
	return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"
}

func (t *Tokenizer) consumeIdentifier() string {
	out := strings.Builder{}
	for {
		if t.isEnd() {
			break
		}
		char := t.input[t.cursor]
		if strings.ContainsRune(t.validIdentifierChars(), rune(char)) {
			out.WriteByte(char)
			t.cursor++
		} else {
			break
		}
	}
	return out.String()
}

func (t *Tokenizer) consumeNumber() (string, error) {
	out := strings.Builder{}
	has_dot := false
loop:
	for {
		if t.isEnd() {
			break
		}

		char := t.input[t.cursor]
		switch {
		case unicode.IsDigit(rune(char)):
			out.WriteByte(char)
			t.cursor++
		case char == '.':
			{
				if !has_dot {
					out.WriteByte(char)
					t.cursor++
				} else {
					return "", fmt.Errorf("consume number: invalid number '%s': multiple decimal points", out.String())
				}
			}
		case char == '_':
			t.cursor++
		default:
			break loop
		}
	}
	return out.String(), nil
}

func (t *Tokenizer) Scan() ([]Token, error) {
	tokens := []Token{}
	tok, err := t.Next()
	if err != nil {
		return nil, fmt.Errorf("scan: %s", err)
	}
	tokens = append(tokens, tok)
	for tok.Type != EOF {
		tok, err = t.Next()
		if err != nil {
			return nil, fmt.Errorf("scan: %s\n", err)
		}
		tokens = append(tokens, tok)
	}
	return tokens, nil
}
