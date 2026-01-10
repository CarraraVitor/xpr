package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const DEBUG = false;

func (exp Number) Eval(env *Env) Type {
	return newFloat(float64(exp))
}

func (s String) Eval(env *Env) Type {
	return newStr(string(s))
}

func (op BinOp) Eval(env *Env) Type {
	right_ := op.right.Eval(env)

	if op.op.Type == EQUAL {
		var_, ok := op.left.(Var)
		if !ok {
			fmt.Printf("ERROR: invalid variable for assignment: '%s'\n", op.left)
			os.Exit(1)
		}
		env.vars[var_] = right_
		return right_
	}

	left_  := op.left.Eval(env)
	if left_.kind != TYPE_FLOAT {
    	fmt.Printf("ERROR: invalid operand for binary operator: '%s'\n", op.left)
		os.Exit(1)
	}

  	if right_.kind != TYPE_FLOAT {
    	fmt.Printf("ERROR: invalid operand for binary operator: '%s'\n", op.left)
		os.Exit(1)
	}

	left := left_.as.float
	right := right_.as.float

	switch op.op.Type {
	case PLUS:
		return newFloat(left + right)
	case MINUS:
		return newFloat(left - right)
	case MULT:
		return newFloat(left * right)
	case DIV:
		return newFloat(left / right)
	case GREATER:
		res := left > right
		if res {
			return newFloat(1.0)
		} else {
			return newFloat(0.0)
		}
	case GREATER_EQUAL: 
		res := left >= right
		if res {
			return newFloat(1.0)
		} else {
			return newFloat(0.0)
		}
	case LESS: 
		res := left < right
		if res {
			return newFloat(1.0)
		} else {
			return newFloat(0.0)
		}
	case LESS_EQUAL: 
		res := left <= right
		if res {
			return newFloat(1.0)
		} else {
			return newFloat(0.0)
		}
	case EQUAL_EQUAL: 
		res := left == right
		if res {
			return newFloat(1.0)
		} else {
			return newFloat(0.0)
		}
	default:
		return newFloat(0.0)
	}
}

func (unop UnOp) Eval(env *Env) Type {
	right := unop.right.Eval(env)
	if right.kind != TYPE_FLOAT {
		fmt.Printf("ERROR: invalid operand for unary operator '%s'\n", unop.right)
		os.Exit(1)
	}
	res := Type{ 
		kind: TYPE_FLOAT,
	}
	switch unop.op.Type {
	case PLUS:
		res.as.float = +right.as.float
		return res
	case MINUS:
		res.as.float = -right.as.float
		return res
	}
	fmt.Printf("ERROR: invalid unary operator: %s\n", unop)
	os.Exit(1)
	return Type{}
}

func (v Var) Eval(env *Env) Type {
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

func (block Block) Eval(env *Env) Type {
	var res Type

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

func (i If) Eval(env *Env) Type {
	cond := i.cond.Eval(env)
	if cond.kind != TYPE_FLOAT {
		fmt.Printf("ERROR: invalid condition in while: '%s'\n", i.cond)
		os.Exit(1)
	}
	if cond.as.float > 0.0 {
		return i.then.Eval(env)
	}
	return Type{
		kind: TYPE_FLOAT,
		as: As{
			float: 0.0,
		},
	}
}

func (ie IfElse) Eval(env *Env) Type {
	cond := ie.cond.Eval(env)
	if cond.kind != TYPE_FLOAT {
		fmt.Printf("ERROR: invalid condition in while: '%s'\n", ie.cond)
		os.Exit(1)
	}
	if cond.as.float > 0.0 {
		return ie.then.Eval(env)
	} else {
		return ie.elze.Eval(env)
	}
}

func (w While) Eval(env *Env) Type {
	res := Type{}
	cond := w.cond.Eval(env)
	if cond.kind != TYPE_FLOAT {
		fmt.Printf("ERROR: invalid condition in while: '%s'\n", w.cond)
		os.Exit(1)
	}
	for cond.as.float > 0.0 {
		res = w.then.Eval(env)
		cond = w.cond.Eval(env)
		if cond.kind != TYPE_FLOAT {
			fmt.Printf("ERROR: invalid condition in while: '%s'\n", w.cond)
			os.Exit(1)
		}
	}
	return res
}

func (p Print) Eval(env *Env) Type {
	res := p.expr.Eval(env)
	fmt.Printf("%s", res)
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
	fmt.Printf("Final Value: %s\n",res) 
}

func REPL(parser *Parser) {
	for {
		scan := bufio.NewScanner(os.Stdin)
		fmt.Printf(">>> ")
		line := ""
		blocks := 0
		scan.Scan()
		txt := scan.Text()
		if strings.Contains(txt, "{")  {
			blocks++
		}
		if strings.Contains(txt, "}") {
			blocks--
		}
		line += txt;
        for blocks > 0 {
			scan.Scan()
			txt := scan.Text()
			if strings.Contains(txt, "{")  {
				blocks++
			}
			if strings.Contains(txt, "}") {
				blocks--
			}
			line += txt;
		}

		tokenizer := NewTokenizer(line)
		tokens, err := tokenizer.Scan()

		if err != nil {
			fmt.Printf("ERROR: Tokenizer: %s\n", err)
			continue
		}

		parser.ResetTokens(tokens)
		expr := parser.Expression(0)
		res := expr.Eval(parser.env)
		fmt.Printf("%s\n",res) 
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
