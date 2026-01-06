package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []Token
	cursor int
}

type Expr interface {
	Eval() float64
	String() string
}

type ExprNumber float64

type BinOp struct {
	left  Expr
	right Expr
	op    Token
}

type UnOp struct {
	right Expr
	op    Token
}

func (exp ExprNumber) String() string {
	return fmt.Sprintf("%.2f", exp)
}


func (binop BinOp) String() string {
	out := strings.Builder{}
	out.Write([]byte{'(', ' '})
	switch binop.op.Type {
	case PLUS:
		out.WriteByte('+')
	case MINUS:
		out.WriteByte('-')
	case MULT:
		out.WriteByte('*')
	case DIV:
		out.WriteByte('/')
	default:
		out.WriteByte('?')
	}
	out.WriteByte(' ')
	out.Write([]byte(binop.left.String()))
	out.WriteByte(' ')
	out.Write([]byte(binop.right.String()))
	out.Write([]byte{' ', ')'})
	return out.String()
}

func (unop UnOp) String() string {
	out := strings.Builder{}
	out.Write([]byte{'(', ' '})
	switch unop.op.Type {
	case PLUS:
		out.WriteByte('+')
	case MINUS:
		out.WriteByte('-')
	}
	out.Write([]byte(unop.right.String()))
	out.Write([]byte{' ', ')'})
	return out.String()
}

func exprNumber(t Token) (ExprNumber, error) {
	if t.Type != NUMBER {
		return 0, fmt.Errorf("expr number: invalid number '%s'", t.Value)
	}
	n, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("expr number: invalid number '%s': failed to parse: strconv: %s", t.Value, err)
	}
	return ExprNumber(n), nil
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
		cursor: 0,
	}
}

func (p *Parser) Peek() Token {
	if p.cursor >= len(p.tokens) {
		return NewEOF()
	}
	tok := p.tokens[p.cursor]
	return tok
}

func (p *Parser) Next() Token {
	if p.cursor >= len(p.tokens) {
		return NewEOF()
	}

	tok := p.tokens[p.cursor]
	p.cursor++
	return tok
}

func (p *Parser) Expect(tok Token) error {
	next := p.Next()
	if next.Type != tok.Type {
		return fmt.Errorf("expected %s, got %s", tok.Type, next.Type)
	}
	return nil
}

func (p *Parser) Parse(min_bp int) Expr {
	var left Expr

	left_tok := p.Next()
	switch left_tok.Type {
	case NUMBER:
		var err error
		left, err = exprNumber(left_tok)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	case LEFT_PAREN:
		left = p.Parse(0)
		err := p.Expect(NewRightParen())
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	case PLUS, MINUS:
		_, rbp := prefixBindingPower(left_tok.Type)
		left = UnOp{
			op:    left_tok,
			right: p.Parse(rbp),
		}
	default:
		return left
	}

	for {
		op := p.Peek()

		lbp, rbp := infixBindingPower(op.Type)
		if lbp <= min_bp {
			break
		}
		p.Next()

		right := p.Parse(rbp)

		if right == nil {
			break
		}

		left = BinOp{
			left:  left,
			right: right,
			op:    op,
		}
	}

	return left
}

func infixBindingPower(toktype TokenType) (int, int) {
	switch toktype {
	case PLUS, MINUS:
		return 1, 2
	case MULT, DIV:
		return 3, 4
	case EOF:
		return 0, 0
	}
	return -1, -1
}

func prefixBindingPower(toktype TokenType) (int, int) {
	switch toktype {
	case PLUS, MINUS:
		return -1, 5
	case EOF:
		return -1, 0
	}
	return -1, -1
}

