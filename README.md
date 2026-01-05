# Xpr

Simple interpreted dynamic language, made for educational purposes.
This project is an attempt in implementing some of Rust's expression-based
features in an interpreted manner.
The main algorithm used for parsing is Pratt's Top-Down Recursive Descent.

Here are some snippets of how the language is supposed to look like when working:
```js
let x(a) = {
    for range 1..100 {
        a *= it;
    }
    a;
};

let y = if x > 100 {
    x * 2;
} else {
    x * 10;
};
```

## References:
- matklad: https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html (https://github.com/matklad/minipratt)
- Robert Nystrom: https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy
- Douglas Crockford: https://www.crockford.com/javascript/tdop/tdop.html
