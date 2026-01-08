package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

const DEBUG = false;

func (exp Number) Eval() float64 {
	return float64(exp)
}

func (op BinOp) Eval() float64 {
	switch op.op.Type {
	case PLUS:
		return op.left.Eval() + op.right.Eval()
	case MINUS:
		return op.left.Eval() - op.right.Eval()
	case MULT:
		return op.left.Eval() * op.right.Eval()
	case DIV:
		return op.left.Eval() / op.right.Eval()
	case EQUAL:
		right := op.right.Eval()
		env[Var(op.left.String())] = right
		return right
	default:
		return 0
	}
}

func (unop UnOp) Eval() float64 {
	switch unop.op.Type {
	case PLUS:
		return +unop.right.Eval()
	case MINUS:
		return -unop.right.Eval()
	}
	fmt.Printf("ERROR: invalid unary operator: %s\n", unop)
	os.Exit(1)
	return 0
}

func (v Var) Eval() float64 {
	val, ok := env[v]
	if !ok {
		return 0
	}
	return val
}

func (block Block) Eval() float64 {
	var res float64
	for _, expr := range block.exprs {
   		res = expr.Eval()
	}
	return res
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
	parser.Parse()
	// fmt.Printf("%f\n", expr.Eval())
	// for parser.ShouldRun()  {
	// 	expr := parser.Parse(0)
	// 	if expr == nil {
	// 		break
	// 	}
	// 	if DEBUG {
	// 		fmt.Printf("---\n")
	// 		fmt.Printf("%s\n", expr)
	// 	}
	// 	fmt.Printf("%f\n", expr.Eval())
	// }
}

func interpret_file(input_file string) {
	expr, err := os.ReadFile(input_file)
	if err != nil {
		fmt.Printf("ERROR: could not read file '%s': %s", input_file, err)
		return
	}
	Interpret(string(expr))
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
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		line := scan.Text()
		Interpret(line)
	}
}
