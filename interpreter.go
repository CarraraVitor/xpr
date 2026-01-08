package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
		fmt.Printf("Evaluationg Expr: %s\n", expr)
   		res = expr.Eval(block.env)
		fmt.Printf("Res: %.2f\n", res)
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

func Interpret(input string) {
	tokenizer := NewTokenizer(input)
	tokens, err := tokenizer.Scan()
	if DEBUG {
		for _, token := range tokens {
			fmt.Printf("%s\n", token)
		}
	}

	if err != nil {
		fmt.Printf("ERROR: Tokenizer: %s\n", err)
		// os.Exit(1)
	}

	parser := NewParser(tokens)
	main_block := parser.Parse()
	res := main_block.Eval(nil)
	// fmt.Printf("%s\n", main_block)
	fmt.Printf("Final Value: %.2f\n",res) 
}

func interpret_file(input_file string) {
	code, err := os.ReadFile(input_file)
	if err != nil {
		fmt.Printf("ERROR: could not read file '%s': %s", input_file, err)
		return
	}
	Interpret(string(code))
}

func main() {
	input := flag.String("input", "", "Input file with source code" )
	flag.Parse()
	if *input != "" {
		interpret_file(*input)
		return
	}
	for {
		input := ""
		fmt.Printf(">>> ")
		if input == "quit" {
			break
		}
		line := ""
		scan := bufio.NewScanner(os.Stdin)
		for !endsWith(line, ';') {
			scan.Scan()
			line += scan.Text()
		}
		Interpret(line)
	}
}

func endsWith(s string, pattern byte) bool {
	if len(s) == 0 {
		return false
	}
	return s[len(s) - 1] == pattern
}
