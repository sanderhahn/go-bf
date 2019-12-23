# Brainfuck

The interpreter could be more efficient if the interpeter did a single pass to determine code jump locations for loops. However the assignment required the use of `io.Reader` and processing instructions without knowing all input at once. [Extended Brainfuck](https://esolangs.org/wiki/Extended_Brainfuck) requires reading the data behind the program file to initialize storage.

```bash
$ go install ./...

$ bf examples/life.bf

$ bf examples/hannoi.bf

$ bf examples/mandelbrot.bf

# compile bf to c
$ bf examples/dbf2c.bf <examples/mandelbrot.bf >examples/mandelbrot.c
$ gcc examples/mandelbrot.c -o examples/mandelbrot
$ ./examples/mandelbrot

# run short test
$ go test -test.short

# test coverage
$ go test -test.short -cover -coverprofile=coverage.out && go tool cover -html=coverage.out
```

- [Brainfuck](http://www.linusakesson.net/programming/brainfuck/index.php)
- [Brainfuck Algorithms](https://esolangs.org/wiki/Brainfuck_algorithms)
- [Brainfuck Examples](http://esoteric.sange.fi/brainfuck/bf-source/prog/)
- [Brainfuck Debugger](http://jsfiddle.net/egon/PyyV2/20/embedded/result/)

## Genetic Programming

Th `bfgen` tool for generating programs was inspired by the research paper
[AI Programmer: Autonomously CreatingSoftware Programs Using Genetic Algorithms](https://arxiv.org/pdf/1709.05703.pdf).
More information is available at [Using Artificial Intelligence to Write Self-Modifying/Improving Programs](http://www.primaryobjects.com/2013/01/27/using-artificial-intelligence-to-write-self-modifying-improving-programs/).

[Genetic programming](https://en.wikipedia.org/wiki/Genetic_programming) uses
genetic evolution in a population of random programs to adapt them into
programs that are increasingly more fit to solve a problem. The fitness of a program
is calculated to see how good the solution is with respect to a certain outcome.
Mutations are introduced in the population by imperfect copying of existing
code. The most fit programs are selected and crossbred to evolve into possibly
better programs.

The brainfuck interpreter is extended to allow more sloppy versions of the code.
The sloppy version will not error on unmached looping operators, so that invalid
loops can be randomly introduced and removed. The `Normalize` function is used
to fix unbalanced brackets so that programs are compatible with more strict
interpreters. Runtime cost for executing a program can be limited and is
returned for evaluation.

## Examples

```bash
$ echo "I Feel Like a Computer" | bfgen
```

```
>+>[]<+>++++[<++++>-]<[<++++>-]<+.-----------------------------------------.++>+
+++[<+++++><++++>-]<.->++++[<++++++++>-]<..+++++++.[-]>++++[<++++++++>-]<.+>++++
+++[<++++++>-]<+.+>++++[<+++++++>-]<.++.------.[-]>++++[<++++><++++>-]<.+>++++++
++[<++++++++>-]<.+>++++[<++++>-]<>>++++<++++++++>++[<++++>-]<.>-++++++[<+++++++>
-]<.++>+-++++++[<+++++++>-]<.--.+++.++++><+.-.---------------.<.[-]++++++++++.
```

<!-- https://www.youtube.com/watch?v=G0-PxhDZV00 -->

The runtime defaults to 10000, but sometimes its beneficial to limit or extend
it fit the length of the text.

```bash
$ cat <<EOF | bfgen -runtime 20000
1
22
333
4444
55555
666666
7777777
88888888
999999999
EOF
```

```
+>++++++[<++++++++>-]<.>++++++++++.>+>+++++++[<+++++++>-]<..>>+<++++++++++.>+>++
+++++[<+++++++>-]<...[+]>+<++++++++++.>+++++[<+++++++>-]<....>++++++++++.<+.>++>
-<<...>+>++++++[<++++++++>-]<.>+<[-]>+<++++++++++.<[+]+>>++>+++++++[<+++++++>-]>
+++++>+++++++[<+++++++>-]<....<<+..[-]++++++++++.<[+]>++++[<++++>-]+++++++<+>+++
>>[-]->+++++++[<++++++++>-]<.>>+<+++>-<<......>+<[-]>+<++++++++++.>++>+++++++[<+
++++++>-]>+++++++>+++++++[<+++++++>-]<......<<..<<<.<......<+++++++>...<++.
```

```bash
# Benchmark add print all byte value representations
$ go test -benchmem -run=^$ github.com/sanderhahn/go-bf -bench "^(BenchmarkAscii)$"
```

```
0x00 = .
0x01 = +.
0x02 = ++.
0x03 = +++.
...
```

## Limitations

There is only one pool so its possible that the population will get stuck in
a solution that doesn't further improve, especially for longer texts. This
situation can be improved by evaluating programs within their own generational
pool. The weight function values early matching letters in output higher and
will start to optimize for program length once a solution is found.

The generator will not generate input `,` instructions because EOF handling is
inconsistent between different implementations.
Actually generating programs that handle input/output in a logical way requires
specifying interaction patterns in a language like [Expect](https://en.wikipedia.org/wiki/Expect).
Otherwise the generator will just use input as source of integer values and
this doesn't result in programs that perform meaningful interactions.
