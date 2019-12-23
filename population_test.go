package bf

import "testing"

import "fmt"

func runPopulation(expected string, maxRuntime, maxIterations int, stopOnSuccess bool) *Population {
	p := NewPopulation()
	p.Expected = expected
	p.MaxRuntime = maxRuntime
	p.MaxManipulation = 1
	for i := 1; i < maxIterations; i++ {
		p.EvaluateAndMutate()
		if stopOnSuccess && p.Fittest().success {
			break
		}
	}
	return p
}

func TestHiPopulation(t *testing.T) {
	hi := runPopulation("hi\n", 200, 100, true)
	_ = hi.Fittest().String()
	if _, ok := hi.SuccessCode(); !ok {
		t.Fail()
	}
}

func TestHelloWorldPopulation(t *testing.T) {
	p := runPopulation("hello world\n", 1000, 500, true)
	if _, ok := p.SuccessCode(); !ok {
		t.Fail()
	}
}

func BenchmarkAscii(b *testing.B) {
	// Finding an optimal solution might require multiple runs
	for i := ' '; i <= '~'; i++ {
		s := string(i)
		p := runPopulation(s, 200, 50, false)
		code, ok := p.SuccessCode()
		if !ok {
			b.FailNow()
		}
		fmt.Printf("%c = %s\n", i, code)
	}
}
