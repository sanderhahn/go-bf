package bf

import "testing"

import "fmt"

func runPopulation(expected []byte, maxRuntime, maxIterations int, stopOnSuccess bool) *Population {
	p := NewPopulation()
	p.Expected = expected
	p.MaxRuntime = maxRuntime
	p.MaxManipulate = 1
	for i := 1; i < maxIterations; i++ {
		p.EvaluateAndMutate()
		if stopOnSuccess && p.Fittest().success {
			break
		}
	}
	return p
}

func TestHiPopulation(t *testing.T) {
	hi := runPopulation([]byte("hi\n"), 200, 100, true)
	_ = hi.Fittest().String()
	if _, ok := hi.SuccessCode(); !ok {
		t.Fail()
	}
}

func TestHelloWorldPopulation(t *testing.T) {
	p := runPopulation([]byte("hello world\n"), 1000, 500, true)
	if _, ok := p.SuccessCode(); !ok {
		t.Fail()
	}
}

func TestHighByte(t *testing.T) {
	p := runPopulation([]byte{byte(0xff)}, 256, 200, false)
	_, ok := p.SuccessCode()
	if !ok {
		t.Fail()
	}
}

func BenchmarkAscii(b *testing.B) {
	// https://esolangs.org/wiki/Brainfuck_constants
	// Finding an optimal solution might require multiple runs
	for i := 0; i <= 255; i++ {
		p := runPopulation([]byte{byte(i)}, 256, 20, false)
		code, ok := p.SuccessCode()
		if !ok {
			b.FailNow()
		}
		fmt.Printf("0x%02x = %s\n", i, code)
	}
}

func TestCharacterFitness(t *testing.T) {
	// byte overflows are considered
	if characterFitness('\x80', '\xff') < characterFitness('\x80', '\x00') {
		t.Fail()
	}
	if characterFitness('\xf2', '\xff') < characterFitness('\xf2', '\x00') {
		t.Fail()
	}
}
