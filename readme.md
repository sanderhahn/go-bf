# Brainfuck

The interpreter could be more efficient if the interpeter did a single pass to determine code jump locations for loops. However the assignment required the use of `io.Reader` and processing instructions without knowing all input at once. [Extended Brainfuck](https://esolangs.org/wiki/Extended_Brainfuck) requires reading the data behind the program file to initialize storage.

```bash
# run short test
$ go test -test.short

$ go install github.com/sanderhahn/bf/...

$ bf examples/life.bf

$ bf examples/hannoi.bf

$ bf examples/mandelbrot.bf

# compile bf to c
$ bf examples/dbf2c.bf <examples/mandelbrot.bf >examples/mandelbrot.c

$ gcc examples/mandelbrot.c -o examples/mandelbrot

$ ./examples/mandelbrot
```

- [Brainfuck](http://www.linusakesson.net/programming/brainfuck/index.php)
- [Brainfuck Algorithms](https://esolangs.org/wiki/Brainfuck_algorithms)
- [Brainfuck Examples](http://esoteric.sange.fi/brainfuck/bf-source/prog/)
- [Conway's Game of Life](http://pi.math.cornell.edu/~lipa/mec/lesson6.html)
