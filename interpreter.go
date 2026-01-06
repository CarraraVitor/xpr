package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func (exp ExprNumber) Eval() float64 {
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

func Interpret(input string) {
	tokenizer := NewTokenizer(input)
	tokens, err := tokenizer.Scan()
	for _, token := range tokens {
		fmt.Printf("%s\n", token)
	}

	if err != nil {
		fmt.Printf("ERROR: Tokenizer: %s\n", err)
		os.Exit(1)
	}

	parser := NewParser(tokens)
	exp := parser.Parse(0)
	fmt.Printf("%s\n", exp)
	fmt.Printf("RES: %f\n", exp.Eval())
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
