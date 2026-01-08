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

type Var string

type Env map[Var]float64
var env = Env{}

type Number float64

type Block struct {
	exprs []Expr
}

type BinOp struct {
	left  Expr
	right Expr
	op    Token
}

type UnOp struct {
	right Expr
	op    Token
}

type Assignment struct {
	left  Expr
	right Expr
}

func (exp Number) String() string {
	return fmt.Sprintf("%.2f", exp)
}

func (binop BinOp) String() string {
	out := strings.Builder{}
	out.Write([]byte("BinOp {\n  op: "))
	switch binop.op.Type {
	case PLUS:
		out.WriteByte('+')
	case MINUS:
		out.WriteByte('-')
	case MULT:
		out.WriteByte('*')
	case DIV:
		out.WriteByte('/')
	case EQUAL:
		out.WriteByte('=')
	default:
		out.WriteByte('?')
	}
	out.Write([]byte("\n  left: "))
	out.Write([]byte(binop.left.String()))
	out.Write([]byte("\n  right: "))
	out.Write([]byte(binop.right.String()))
	out.Write([]byte("\n}"))
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

func (v Var) String() string {
	return string(v)
}

func (block Block) String() string {
	out := strings.Builder{}
	out.Write([]byte("Block: "))
	out.WriteByte('{')
	out.WriteByte('\n')
	for _, expr := range block.exprs {
		out.Write([]byte(expr.String()))
		out.WriteByte('\n')
	}
	out.WriteByte('\n')
	out.WriteByte('}')
	return out.String()
}

func exprNumber(t Token) (Number, error) {
	if t.Type != NUMBER {
		return 0, fmt.Errorf("expr number: invalid number '%s'", t.Value)
	}
	n, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("expr number: invalid number '%s': failed to parse: strconv: %s", t.Value, err)
	}
	return Number(n), nil
}

func exprVar(t Token) (Var, error) {
    if t.Type != ID {
		return Var(""), fmt.Errorf("expr var: invalid identifier'%s'", t.Value)
	}

	variable := Var(t.Value)
	return variable, nil
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: tokens,
		cursor: 0,
	}
}

func (p *Parser) shouldRun() bool {
	return p.cursor < len(p.tokens) 
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

func (p *Parser) Assert(tok Token, typ TokenType) error {
	if tok.Type != typ {
		return fmt.Errorf("Parser Assert: expected '%s', got '%s'", typ, tok.Type)
	}
	return nil
}

func (p *Parser) Parse() {
	b := p.Block()
	fmt.Printf("block: %s\n", b)
	fmt.Printf("RES: %.2f\n", b.Eval())
}

func (p *Parser) Block() Block {
	block := Block{}
	for p.Peek().Type != RIGHT_CURLY {
		expr := p.Expression(0)
		if expr == nil {
			break
		}
		block.exprs = append(block.exprs, expr)
	}
	return block
}

func (p *Parser) Expression(prev_bp int) Expr {
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
		left = p.Expression(0)
		err := p.Expect(NewRightParen())
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	case LEFT_CURLY:
		left = p.Block()
		err := p.Expect(NewRightCurly())
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	case PLUS, MINUS:
		_, rbp := prefixBindingPower(left_tok.Type)
		left = UnOp{
			op:    left_tok,
			right: p.Expression(rbp),
		}
	case ID:
		var err error
		left, err = exprVar(left_tok)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	case SEMICOLON: 
		fmt.Printf("SEMICOLON FOUND!!!!\n")
		return nil
	default:
		return nil
	}

	for {
		op := p.Peek()
		if op.Type == SEMICOLON {
			p.Next()
			return left
		}

		lbp, rbp := infixBindingPower(op.Type)
		if lbp <= prev_bp {
			return left
		}
		p.Next()

		right := p.Expression(rbp)

		if right == nil {
			return left
		}

		left = BinOp{
			left:  left,
			right: right,
			op:    op,
		}
	}
}

func infixBindingPower(toktype TokenType) (int, int) {
	switch toktype {
	case EQUAL:
		return 1, 2
	case PLUS, MINUS:
		return 3, 4
	case MULT, DIV:
		return 5, 6
	case EOF:
		return 0, 0
	}
	return -1, -1
}

// prefixBindingPower returns all the lbp as -1
// because they are meaningless,
// as a prefix operator only operates on operands to its right.
// Besides that, a value is returned for the sake of keeping the
// usage equal to the infixBindingPower function
//
// The expected usage is something like:
//
//	_, rbp := prefixBindingPower(toktype)
func prefixBindingPower(toktype TokenType) (int, int) {
	switch toktype {
	case PLUS, MINUS:
		return -1, 8
	case EOF:
		return -1, 0
	}
	return -1, -1
}
