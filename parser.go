package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)


type Parser struct {
	tokens    []Token
	cursor    int
	env       *Env
}

type Expr interface {
	Eval(env *Env) float64
	String() string
}

type Env struct {
	vars map[Var]float64
	parent *Env
}

type Number float64

type BinOp struct {
	left  Expr
	right Expr
	op    Token
}

type UnOp struct {
	right Expr
	op    Token
}

type Var string

type Block struct {
	exprs   []Expr
	env     *Env
}

type Assignment struct {
	left  Expr
	right Expr
}

type If struct {
	cond Expr
	then Block
}

type IfElse struct {
	If
	elze Block
}

type While struct {
	cond Expr
	then Block
}

type Print struct {
	expr Expr
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
	case GREATER:
		out.WriteByte('>')
	case GREATER_EQUAL: 
		out.WriteByte('>')
		out.WriteByte('=')
	case LESS: 
		out.WriteByte('<')
	case LESS_EQUAL: 
		out.WriteByte('<')
		out.WriteByte('=')
	case EQUAL_EQUAL: 
		out.WriteByte('=')
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

func (i If) String() string {
	return fmt.Sprintf("if (%s) {\n%s\n}\n", i.cond, i.then)
}

func (ie IfElse) String() string {
	return fmt.Sprintf("if (%s) {\n%s\n} else {\n%s\n}", ie.cond, ie.then, ie.elze)
}

func (w While) String() string {
	return fmt.Sprintf("while (%s) {\n%s\n}\n", w.cond, w.then)
}

func (p Print) String() string {
	return fmt.Sprintf("PRINT: %s", p.expr.String())
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

func newEnv() Env {
	return Env{
		vars: make(map[Var]float64),
		parent: nil,
	}
}

func NewParser(tokens []Token) Parser {
	env := newEnv()
	return Parser{
		tokens: tokens,
		cursor: 0,
		env: &env,
	}
}

func (p *Parser) ResetTokens(tokens []Token) {
	p.tokens = tokens
	p.cursor = 0
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

func (p *Parser) Parse() Expr {
	return p.Block()
}

func (p *Parser) IfElse() (res Expr, err error) {
	cond := p.Expression(0)
	err = p.Expect(NewLeftCurly())
	if err != nil { return }
	then := p.Block()
	err = p.Expect(NewRightCurly())
	if err != nil { return }

	res = If{
		cond: cond,
		then: then,
	}

	if p.Peek().Type == ELSE {
		p.Next()
		err = p.Expect(NewLeftCurly())
		if err != nil { return }

		elze := p.Block()

		err = p.Expect(NewRightCurly())
		if err != nil { return }

		res = IfElse{
			If:   res.(If),
			elze: elze,
		}
	} 

	return
}

func (p *Parser) While() (w Expr, err error) {
	w = While{}

	cond := p.Expression(0)
	err = p.Expect(NewLeftCurly())
	if err != nil { return }
	then := p.Block(p.env.vars)
	err = p.Expect(NewRightCurly())
	if err != nil { return }

	w = While{
		cond: cond,
		then: then,
	}
	return
}

func (p *Parser) Block(vars_ ...map[Var]float64) Block {
	var vars map[Var]float64
	if len(vars_) > 0 {
		vars = vars_[0]
	} else {
		vars = make(map[Var]float64)
	}
	new_env := Env{
		vars: vars,
		parent: p.env,
	}
	block := Block{ env: &new_env }
	p.env = block.env
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
	case IF:
		var err error
		left, err = p.IfElse()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	case WHILE:
		var err error
		left, err = p.While()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	case PRINT:
		left = Print{
			expr: p.Expression(0),
		}
	case SEMICOLON: 
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
	case GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, EQUAL_EQUAL:
		return 3, 4
	case PLUS, MINUS:
		return 5, 6
	case MULT, DIV:
		return 7, 8
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
		return -1, 10
	case EOF:
		return -1, 0
	}
	return -1, -1
}
