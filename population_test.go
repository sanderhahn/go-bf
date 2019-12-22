package bf

import "testing"

func RunPopulation(expected string, maxRuntime, maxIterations int) *Population {
	p := NewPopulation()
	p.Expected = expected
	p.MaxRuntime = maxRuntime
	for i := 1; i < maxIterations; i++ {
		p.EvaluateAndMutate()
		if p.Fittest().success {
			break
		}
	}
	return p
}

func TestHiPopulation(t *testing.T) {
	hi := RunPopulation("hi\n", 200, 100)
	_ = hi.Fittest().String()
	if _, ok := hi.SuccessCode(); !ok {
		t.Fail()
	}
}

func TestHelloWorldPopulation(t *testing.T) {
	p := RunPopulation("hello world\n", 1000, 500)
	if _, ok := p.SuccessCode(); !ok {
		t.Fail()
	}
}
