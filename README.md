# Xpr

Simple interpreted dynamic language, made for educational purposes.
This project is an attempt in implementing some of Rust's expression-based
features in an interpreted manner.
The main algorithm used for parsing is Pratt's Top-Down Recursive Descent.

# Examples
[Imperative Fibonacci](./examples/fibonacci.xpr)
[Recursive Fibonacci](./examples/rec_fibonacci.xpr)

# Usage:
For now there is no releases, so you need the Go compiler.
```sh
git clone https://github.com/CarraraVitor/xpr .
go run . -input ./examples/fibonacci.xpr
```


## References:
- matklad: https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html (https://github.com/matklad/minipratt)
- Robert Nystrom: https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy
- Douglas Crockford: https://www.crockford.com/javascript/tdop/tdop.html
