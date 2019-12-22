package bf

import (
	"bytes"
	"errors"
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

func TestStrict(t *testing.T) {
	i := NewInterpreter(nil, nil)
	if err := i.Interpret(code(`[`)); err == nil {
		t.Fail()
	}
	if err := i.Interpret(code(`]`)); err == nil {
		t.Fail()
	}
	if err := i.Interpret(code(`+[`)); err == nil {
		t.Fail()
	}
}

func TestSloppy(t *testing.T) {
	out := &strings.Builder{}
	i := NewInterpreter(out, nil)
	_, err := i.InterpretExtended(code(`[`), false, -1)
	if err != nil {
		t.Fail()
	}
	_, err = i.InterpretExtended(code(`]`), false, -1)
	if err != nil {
		t.Fail()
	}
}

func TestRuntime(t *testing.T) {
	out := &strings.Builder{}
	i := NewInterpreter(out, nil)
	_, err := i.InterpretExtended(code(`+++.`), true, 1)
	if err != errExhaustedRuntime {
		t.Fail()
	}
	i = NewInterpreter(out, nil)
	_, err = i.InterpretExtended(code(`+++.`), true, 4)
	if out.String() != "\003" {
		t.Fail()
	}
}

func calc(memory []byte, program string, expected []byte) error {
	i := NewInterpreter(nil, nil)
	copy(i.memory, memory)
	err := i.Interpret(code(program))
	if err != nil {
		return err
	}
	if i.ptr != 0 || !bytes.Equal(i.memory[:len(expected)], expected) {
		return errors.New("wrong calc")
	}
	return nil
}

func TestCalc(t *testing.T) {
	if err := calc([]byte{255}, `[-]`, []byte{0}); err != nil {
		t.Fail()
	}
	if err := calc([]byte{}, `>++++[<++++>-]<`, []byte{16}); err != nil {
		t.Fail()
	}
	if err := calc([]byte{}, `>++++[<++++++++>-]<`, []byte{32}); err != nil {
		t.Fail()
	}
	if err := calc([]byte{}, `>++++++[<++++++++>-]<`, []byte{48}); err != nil {
		t.Fail()
	}
	if err := calc([]byte{}, `>++++++++[<++++++++>-]<`, []byte{64}); err != nil {
		t.Fail()
	}
	// dup
	if err := calc([]byte{4}, `[>+>+<<-]>>[<<+>>-]<<`, []byte{4, 4, 0}); err != nil {
		t.Fail()
	}
	// sum
	if err := calc([]byte{2, 2}, `>[<+>-]<`, []byte{4, 0}); err != nil {
		t.Fail()
	}
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
