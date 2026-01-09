package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const DEBUG = false;

func (exp Number) Eval(env *Env) float64 {
	return float64(exp)
}

func (op BinOp) Eval(env *Env) float64 {
	switch op.op.Type {
	case PLUS:
		return op.left.Eval(env) + op.right.Eval(env)
	case MINUS:
		return op.left.Eval(env) - op.right.Eval(env)
	case MULT:
		return op.left.Eval(env) * op.right.Eval(env)
	case DIV:
		return op.left.Eval(env) / op.right.Eval(env)
	case EQUAL:
		right := op.right.Eval(env)
		env.vars[Var(op.left.String())] = right
		return right
	case GREATER:
		res := op.left.Eval(env) > op.right.Eval(env)
		if res {
			return 1.0
		} else {
			return 0.0
		}
	case GREATER_EQUAL: 
		res := op.left.Eval(env) >= op.right.Eval(env)
		if res {
			return 1.0
		} else {
			return 0.0
		}
	case LESS: 
		res := op.left.Eval(env) < op.right.Eval(env)
		if res {
			return 1.0
		} else {
			return 0.0
		}
	case LESS_EQUAL: 
		res := op.left.Eval(env) <= op.right.Eval(env)
		if res {
			return 1.0
		} else {
			return 0.0
		}
	case EQUAL_EQUAL: 
		res := op.left.Eval(env) == op.right.Eval(env)
		if res {
			return 1.0
		} else {
			return 0.0
		}
	default:
		return 0
	}
}

func (unop UnOp) Eval(env *Env) float64 {
	switch unop.op.Type {
	case PLUS:
		return +unop.right.Eval(env)
	case MINUS:
		return -unop.right.Eval(env)
	}
	fmt.Printf("ERROR: invalid unary operator: %s\n", unop)
	os.Exit(1)
	return 0
}

func (v Var) Eval(env *Env) float64 {
	if env == nil {
		fmt.Printf("ERROR: unknown variable '%s'\n", v)
		os.Exit(1)
	}

	if DEBUG {
		fmt.Printf("----\n")
		fmt.Printf("Evaluating Var '%s'\n", v)
		fmt.Printf("Env Addr: %p\n", env)
		fmt.Printf("Env: %#v\n", env)
	}

	val, ok := env.vars[v]
	if !ok {
		return v.Eval(env.parent)
	}
	return val
}

func (block Block) Eval(env *Env) float64 {
	var res float64

	if DEBUG {
		fmt.Printf("-----\n")
		fmt.Printf("Block:\n%s\n", block)
		fmt.Printf("Env:\n%#v\n", block.env)
	}

	for _, expr := range block.exprs {
		// fmt.Printf("Evaluating Expr: %s\n", expr)
		// fmt.Printf("Res: %.2f\n", res)
   		res = expr.Eval(block.env)
	}
	return res
}

func (i If) Eval(env *Env) float64 {
	cond := i.cond.Eval(env)
	if cond > 0 {
		return i.then.Eval(env)
	}
	return 0
}

func (ie IfElse) Eval(env *Env) float64 {
	cond := ie.cond.Eval(env)
	if cond > 0 {
		return ie.then.Eval(env)
	} else {
		return ie.elze.Eval(env)
	}
}

func (w While) Eval(env *Env) float64 {
	res := 0.0
	for w.cond.Eval(env) > 0.0 {
		res = w.then.Eval(env)
	}
	return res
}

func (p Print) Eval(env *Env) float64 {
	res := p.expr.Eval(env)
	fmt.Printf("%.2f\n", res)
	return res
}

func interpret_file(parser *Parser, input_file string) {
	code, err := os.ReadFile(input_file)
	if err != nil {
		fmt.Printf("ERROR: could not read file '%s': %s", input_file, err)
		return
	}

	tokenizer := NewTokenizer(string(code))
	tokens, err := tokenizer.Scan()
	if DEBUG {
		for _, token := range tokens {
			fmt.Printf("%s\n", token)
		}
	}

	if err != nil {
		fmt.Printf("ERROR: Tokenizer: %s\n", err)
		return
		// os.Exit(1)
	}

	parser.ResetTokens(tokens)
	main_block := parser.Parse()
	res := main_block.Eval(nil)
	// fmt.Printf("%s\n", main_block)
	fmt.Printf("Final Value: %.2f\n",res) 
}

func scanLine(scan bufio.Scanner) string {
	txt := ""
	for !endsWith(txt, ';') {
		scan.Scan()
		txt += scan.Text()
	}
	return txt
}

func scanBlock(scan bufio.Scanner) string {
	txt := ""
	for !strings.Contains(txt, "}"){
		scan.Scan()
		txt += scan.Text()
		if strings.Contains(txt, "{") && !strings.Contains(txt, "}") {
			txt += scanBlock(scan)
		}
	}
	return txt
}

func REPL(parser *Parser) {
	for {

		fmt.Printf(">>> ")
		line := ""
		scan := bufio.NewScanner(os.Stdin)
		txt := scanLine(*scan)
		if strings.Contains(txt, "{") && !strings.Contains(txt, "}") {
			txt += scanBlock(*scan)
		}
		line += txt;
		fmt.Printf("LINE: %s\n", line)

		tokenizer := NewTokenizer(line)
		tokens, err := tokenizer.Scan()

		if err != nil {
			fmt.Printf("ERROR: Tokenizer: %s\n", err)
			continue
		}

		parser.ResetTokens(tokens)
		expr := parser.Expression(0)
		res := expr.Eval(parser.env)
		fmt.Printf("Final Value: %.2f\n",res) 
	}
}

func main() {
	input := flag.String("input", "", "Input file with source code" )
	flag.Parse()
	tokens := []Token{}
	parser := NewParser(tokens)

	if *input != "" {
		interpret_file(&parser, *input)
		return
	}

	REPL(&parser)
}

func endsWith(s string, pattern byte) bool {
	if len(s) == 0 {
		return false
	}
	return s[len(s) - 1] == pattern
}
