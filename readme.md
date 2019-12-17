# Brainfuck

The interpreter could be more efficient if the interpeter did a single pass to determine code jump locations for loops. However the assignment required the use of `io.Reader` and processing instructions without knowing all input at once.

```bash
# run short test
$ go test -test.short

$ go install github.com/sanderhahn/bf/...

$ bf life.bf
```

- [Brainfuck](http://www.linusakesson.net/programming/brainfuck/index.php)
- [Conway's Game of Life](http://pi.math.cornell.edu/~lipa/mec/lesson6.html)

[Extended Brainfuck](https://esolangs.org/wiki/Extended_Brainfuck) requires reading the data behind the program file to initialize storage.
