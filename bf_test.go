package bf

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func code(code string) io.Reader {
	return bytes.NewReader([]byte(code))
}

func TestHello(t *testing.T) {
	out := &strings.Builder{}
	i := NewInterpreter(out, nil)
	i.Interpret(code(`
	++++++++++
	[>+++++++>++++++++++>+++>+<<<<-] De initialiserende loop om de array te maken
	>++. Print 'H'
	>+. Print 'e'
	+++++++. Print 'l'
	. Print 'l'
	+++. Print 'o'
	>++. Print ' '
	<<+++++++++++++++. Print 'W'
	>. Print 'o'
	+++. Print 'r'
	------. Print 'l'
	--------. Print 'd'
	>+. Print '!'
	>. Print newline
	`))
	expected := "Hello World!\n"
	if out.String() != expected {
		t.Fatalf("%s != %s", out.String(), expected)
	}
}

// go test -bench=. -run=Hannoi
func BenchmarkHannoi(b *testing.B) {
	hannoi, err := ioutil.ReadFile("examples/hannoi.bf")
	if err != nil {
		b.Fail()
	}
	out := &strings.Builder{}
	i := NewInterpreter(out, nil)
	i.Interpret(bytes.NewReader(hannoi))
}

// http://www.linusakesson.net/programming/brainfuck/output.txt
func TestLife(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping life test in short mode")
	}
	life, err := ioutil.ReadFile("examples/life.bf")
	if err != nil {
		t.Fail()
	}
	instructions := `
cc
cd
ce
dc
dd
de
ec
ed
ee
ff
fg
fh
gf
gg
gh
hf
hg
hh








q
`
	out := &strings.Builder{}
	i := NewInterpreter(out, strings.NewReader(instructions))
	i.Interpret(bytes.NewReader(life))
	expected := ` abcdefghij
a----------
b----------
c----------
d----------
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c----------
d----------
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--*-------
d----------
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--**------
d----------
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d----------
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--*-------
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--**------
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e----------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--*-------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--**------
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f----------
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----*----
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----**---
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g----------
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g-----*----
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g-----**---
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g-----***--
h----------
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g-----***--
h-----*----
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g-----***--
h-----**---
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g-----***--
h-----***--
i----------
j----------
> abcdefghij
a----------
b---*------
c--*-*-----
d-*---*----
e--*---*---
f---*---*--
g----*---*-
h-----*-*--
i------*---
j----------
> abcdefghij
a----------
b---*------
c--***-----
d-***-*----
e--*---*---
f---*---*--
g----*-***-
h-----***--
i------*---
j----------
> abcdefghij
a----------
b--***-----
c-*--------
d-*---*----
e-*--*-*---
f---*-*--*-
g----*---*-
h--------*-
i-----***--
j----------
> abcdefghij
a---*------
b--**------
c-*-**-----
d***--*----
e--*-*-*---
f---*-*-*--
g----*--***
h-----**-*-
i------**--
j------*---
> abcdefghij
a--**------
b----------
c*---*-----
d*----*----
e--*-*-*---
f---*-*-*--
g----*----*
h-----*---*
i----------
j------**--
> abcdefghij
a----------
b---*------
c----------
d-*-***----
e---**-*---
f---*-**---
g----***-*-
h----------
i------*---
j----------
> abcdefghij
a----------
b----------
c--**------
d--**-*----
e------*---
f---*------
g----*-**--
h------**--
i----------
j----------
> abcdefghij
a----------
b----------
c--***-----
d--***-----
e--***-----
f-----***--
g-----***--
h-----***--
i----------
j----------
>`
	if out.String() != expected {
		t.Fatalf("%s != %s", out.String(), expected)
	}
}
