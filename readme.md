# BrainFuck

The interpreter could be more efficient [File.Seek](https://golang.org/pkg/os/#File.Seek)
was used. However the assignment required the use of `io.Reader`.

```bash
$ go install github.com/sanderhahn/bf/...

$ bf life.bf
```

- [Brainfuck](http://www.linusakesson.net/programming/brainfuck/index.php)
- [Conway's Game of Life](http://pi.math.cornell.edu/~lipa/mec/lesson6.html)
