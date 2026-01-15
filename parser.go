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
	env    *Env
}

type TypeKind int

const (
	TYPE_FLOAT TypeKind = iota
	TYPE_STRING
)

type As struct {
	float float64
	str   string
}

type Type struct {
	kind TypeKind
	as   As
}

type Expr interface {
	Eval(env *Env) Type
	String() string
}

type Env struct {
	vars   map[Var]Type
	funcs  map[string]*Function
	parent *Env
}

type Number float64

type String string

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
	exprs []Expr
	env   *Env
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

type Function struct {
	name   string
	params []Var
	body   Block
}

type FunctionCall struct {
	fun  *Function
	args map[Var]Expr
}

type Return struct {
	expr Expr
}

func newFloat(f float64) Type {
	return Type{
		kind: TYPE_FLOAT,
		as:   As{float: f},
	}
}

func newStr(s string) Type {
	return Type{
		kind: TYPE_STRING,
		as:   As{str: s},
	}
}

func (env *Env) getFunc(expr Expr) (*Function, error) {
	var_name, ok := expr.(Var)
	if !ok {
		return nil, fmt.Errorf("env get func: invalid function name '%v'", expr)
	}
	name := string(var_name)

	fun, ok := env.funcs[name]
	if !ok {
		if env.parent == nil {
			return nil, fmt.Errorf("env get func: unknown function name: '%s'", name)
		}

		return env.parent.getFunc(expr)
	}
	return fun, nil
}

func (typ Type) String() string {
	var res string
	switch typ.kind {
	case TYPE_FLOAT:
		res = fmt.Sprintf("%.2f", typ.as.float)
	case TYPE_STRING:
		res = typ.as.str
	default:
		res = "?????"
	}
	return res
}

func (exp Number) String() string {
	return fmt.Sprintf("%.2f", exp)
}

func (s String) String() string {
	return string(s)
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

func (fc FunctionCall) String() string {
	out := ""
	name := fc.fun.name
	out += name + "("
	for _, param := range fc.fun.params {
		val, ok := fc.args[param]
		if ok {
			out += fmt.Sprintf("%s = %s,", string(param), val.String())
		} else {
			out += string(param) + ","
		}
	}
	out += ")"
	return out
}

func (p Print) String() string {
	return fmt.Sprintf("PRINT: %s", p.expr.String())
}

func (r Return) String() string {
	return fmt.Sprintf("RETURN: %s", r.expr.String())
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

func exprStr(s Token) (String, error) {
	if s.Type != STR_LIT {
		return String(""), fmt.Errorf("expr var: invalid string literal '%s'", s.Value)
	}
	return String(s.Value), nil
}

func exprVar(t Token) (Var, error) {
	if t.Type != ID {
		return Var(""), fmt.Errorf("expr var: invalid identifier'%s'", t.Value)
	}

	variable := Var(t.Value)
	return variable, nil
}

func exprFunc(name string, params []Var, body Block) Function {
	return Function{
		name:   name,
		params: params,
		body:   body,
	}
}

func newEnv() Env {
	return Env{
		vars:   make(map[Var]Type),
		funcs:  make(map[string]*Function),
		parent: nil,
	}
}

func NewParser(tokens []Token) Parser {
	env := newEnv()
	return Parser{
		tokens: tokens,
		cursor: 0,
		env:    &env,
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
	if err != nil {
		return
	}
	then := p.Block()
	err = p.Expect(NewRightCurly())
	if err != nil {
		return
	}

	res = If{
		cond: cond,
		then: then,
	}

	if p.Peek().Type == ELSE {
		p.Next()
		err = p.Expect(NewLeftCurly())
		if err != nil {
			return
		}

		elze := p.Block()

		err = p.Expect(NewRightCurly())
		if err != nil {
			return
		}

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
	if err != nil {
		return
	}

	env := newEnv()
	env.vars = p.env.vars
	then := p.Block(env)
	err = p.Expect(NewRightCurly())
	if err != nil {
		return
	}

	w = While{
		cond: cond,
		then: then,
	}
	return
}

func (p *Parser) Block(envs ...Env) Block {
	var new_env Env
	if len(envs) > 0 {
		new_env = envs[0]
	} else {
		new_env = newEnv()
	}

	new_env.parent = p.env
	block := Block{env: &new_env}
	p.env = block.env

	for p.Peek().Type != RIGHT_CURLY {
		if p.cursor >= len(p.tokens) {
			break
		}
		expr := p.Expression(0)
		if expr != nil {
			block.exprs = append(block.exprs, expr)
		}
	}

	p.env = block.env.parent
	return block
}

func (p *Parser) FunctionDeclaration() error {
	func_name_tok := p.Next()
	err := p.Assert(func_name_tok, ID)
	if err != nil {
		return fmt.Errorf("function declaration: invalid function name: %s", err)
	}
	name := func_name_tok.Value

	err = p.Expect(NewLeftParen())
	if err != nil {
		return fmt.Errorf("function declaration: expected '(' after function name")
	}

	params := []Var{}
params_loop:
	for {
		typ := p.Peek().Type
		switch typ {
		case ID:
			param, err := exprVar(p.Next())
			if err != nil {
				return fmt.Errorf("function declaration: invalid function parameter: %s", err)
			}
			params = append(params, param)
		case COMMA:
			p.Next()
		default:
			break params_loop
		}
	}

	err = p.Expect(NewRightParen())
	if err != nil {
		return fmt.Errorf("function declaration: expected ')' after function's params")
	}

	err = p.Expect(NewLeftCurly())
	if err != nil {
		return fmt.Errorf("function declaration: expected '{' after function's parameters ")
	}

	fun := &Function{name: name, params: params}
	p.env.funcs[name] = fun

	env := newEnv()
	// for _, param := range params {
	// 	env.vars[param] = newFloat(0.0)
	// }

	body := p.Block(env)
	*fun = exprFunc(name, params, body)

	err = p.Expect(NewRightCurly())
	if err != nil {
		return fmt.Errorf("function declaration: expected '}' after function's body")
	}

	if p.Peek().Type == SEMICOLON {
		p.Next()
	}

	return nil
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
	case STR_LIT:
		var err error
		left, err = exprStr(left_tok)
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
	case FUNCTION:
		err := p.FunctionDeclaration()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		return nil

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
	case RETURN:
		left = Return{
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

		lbp, _ := postfixBindingPower(op.Type)
		if lbp != -1 {
			if lbp <= prev_bp {
				return left
			}

			function, err := p.env.getFunc(left)
			if err != nil {
				fmt.Printf("ERROR: invalid function call: %s\n", err)
				os.Exit(1)
			}

			// NOTE: the only possible postfix operator for now is a(b)
			err = p.Assert(p.Next(), LEFT_PAREN)
			if err != nil {
				fmt.Printf("ERROR: only supported postfix operator is Function Call. Invalid operator: '%s'\n", op.Value)
				os.Exit(1)
			}

			args := []Expr{}

			for peek := p.Peek().Type; peek != RIGHT_PAREN && peek != EOF; peek = p.Peek().Type {
				arg := p.Expression(0)
				args = append(args, arg)
				if p.Peek().Type == COMMA {
					p.Next()
				} else {
					break
				}
			}

			if len(args) > len(function.params) {
				fmt.Printf("ERROR: expected at most %d arguments, but got %d\n", len(function.params), len(args))
				os.Exit(1)
			}

			args_map := make(map[Var]Expr)
			for i := range len(args) {
				arg := args[i]
				param := function.params[i]
				args_map[param] = arg
			}

			err = p.Expect(NewRightParen())
			if err != nil {
				fmt.Printf("ERROR: expected ')' in function call: %s\n", err)
				os.Exit(1)
			}

			left = FunctionCall{
				fun:  function,
				args: args_map,
			}
			continue
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

func postfixBindingPower(toktype TokenType) (int, int) {
	switch toktype {
	case LEFT_PAREN:
		return 9, -1
	}
	return -1, -1
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
