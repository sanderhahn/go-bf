package bf

import (
	"bytes"
	"testing"
)

func TestNormalize(t *testing.T) {
	if !bytes.Equal(Normalize(Program(`[`)), Program(``)) {
		t.Fail()
	}
	if !bytes.Equal(Normalize(Program(`]`)), Program(``)) {
		t.Fail()
	}
	if !bytes.Equal(Normalize(Program(`][`)), Program(``)) {
		t.Fail()
	}
	if !bytes.Equal(Normalize(Program(`[]`)), Program(`[]`)) {
		t.Fail()
	}
}

func TestInsertAt(t *testing.T) {
	if !bytes.Equal(insertAt(Program(`a`), 0, Program(`b`)), Program(`ba`)) {
		t.Fail()
	}
	if !bytes.Equal(insertAt(Program(`a`), 1, Program(`b`)), Program(`ab`)) {
		t.Fail()
	}
}

func TestRemoveAt(t *testing.T) {
	if !bytes.Equal(removeAt(Program(`ab`), 0, 0), Program(`ab`)) {
		t.Fail()
	}
	if !bytes.Equal(removeAt(Program(`ab`), 0, 1), Program(`b`)) {
		t.Fail()
	}
	if !bytes.Equal(removeAt(Program(`ab`), 1, 1), Program(`a`)) {
		t.Fail()
	}
}

func TestReplaceAt(t *testing.T) {
	if !bytes.Equal(replaceAt(Program(`ab`), 0, 1, Program(`c`)), Program(`cb`)) {
		t.Fail()
	}
	if !bytes.Equal(replaceAt(Program(`ab`), 1, 1, Program(`c`)), Program(`ac`)) {
		t.Fail()
	}
}

func TestMutateAt(t *testing.T) {
	// run some mutations to inflate coverage
	for i := 0; i < 50; i++ {
		p := NewRandomProgram(2)
		Mutate(p, []Entry{Entry{program: Program(`x`)}})
	}
	Mutate(Program(``), []Entry{})
}
