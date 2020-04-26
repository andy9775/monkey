import time


def fib(x):
    if x <= 1:
        return x
    return fib(x-1) + fib(x-2)


start = time.time()
for i in range(15):
    print(fib(30))

end = time.time()
print((end - start)/i)  # in seconds
