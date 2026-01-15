def fib(n):
    if n == 0:
        return 0
    if n == 1:
        return 1
    return fib(n-1) + fib(n-2)

def test(n):
    if n < 0.1:
        return 0
    else:
        return n + test(n-1)

print(fib(10))
