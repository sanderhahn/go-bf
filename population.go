package bf

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

const keepSize = 32
const manipulationSize = 32
const populationSize = keepSize * manipulationSize

// Population maintains the pool of programs in the form of entries
type Population struct {
	entries       []Entry
	Expected      []byte
	MaxRuntime    int
	MaxManipulate int
}

// Entry maintains information of a program
type Entry struct {
	program    Program
	output     []byte
	runtime    int
	err        error
	fitness    float64
	success    bool
	generation int
}

// NewPopulation constructor
func NewPopulation() *Population {
	entries := make([]Entry, populationSize)
	for i := range entries {
		entries[i].program = NewRandomProgram(1)
	}
	return &Population{
		entries:       entries,
		MaxRuntime:    10000,
		MaxManipulate: 1,
	}
}

// Fittest is the top performing program
func (p *Population) Fittest() *Entry {
	return &p.entries[0]
}

// SuccessCode returns the program code when success is reached
func (p *Population) SuccessCode() (Program, bool) {
	if p.Fittest().success {
		return Normalize(p.Fittest().program), true
	}
	return nil, false
}

// EvaluateAndMutate will execute programs and evaluate fitness and mutate
func (p *Population) EvaluateAndMutate() {
	for i := range p.entries {
		entry := &p.entries[i]
		entry.fitness = 0.0
		entry.err = entry.exec("", p.MaxRuntime)
		entry.success = entry.err == nil && bytes.Equal(entry.output, p.Expected)
		if entry.err == nil {
			entry.calculateFitness(p.Expected)
		}
	}

	sort.Stable(byFitness(p.entries))

	for i := 0; i < keepSize; i++ {
		keepEntry := &p.entries[i]
		keepEntry.generation++
		for m := 1; m < manipulationSize; m++ {
			entry := &p.entries[(m*keepSize)+i]
			if m == 1 {
				// new generation
				entry.program = NewRandomProgram(1)
				entry.generation = 0
			} else {
				entry.program = append(Program{}, keepEntry.program...)
				manipulation := rand.Intn(p.MaxManipulate)
				if manipulation == 0 {
					manipulation++
				}
				entry.program = Mutate(entry.program, manipulation, p.entries[0:keepSize])
				entry.generation = keepEntry.generation
			}
		}
	}
}

func (e *Entry) exec(input string, maxRuntime int) error {
	output := bytes.NewBuffer([]byte{})
	i := NewInterpreter(output, bytes.NewReader([]byte(input)))
	runtime, err := i.InterpretExtended(bytes.NewReader(e.program), false, maxRuntime)
	e.err = err
	if err == nil {
		e.runtime = maxRuntime - runtime
		e.output = output.Bytes()
	} else {
		e.runtime = 0
		e.output = []byte{}
	}
	return err
}

func characterFitness(expected, actual byte) float64 {
	diff := math.Abs(float64(expected)-float64(actual)) / 255.0
	return 1.0 - diff
}

func (e *Entry) calculateFitness(expected []byte) {
	fitness := 0.0
	factor := 1.0
	ok := true
	for i, ch := range expected {
		if i < len(e.output) {
			ok = ok && ch == e.output[i]
			if ok {
				fitness += 1.0
			} else {
				fitness += (characterFitness(ch, e.output[i]) / factor)
				factor *= 10
			}
		} else {
			ok = false
		}
	}

	lenOutput := float64(len(e.output))
	lenExpected := float64(len(expected))
	if lenOutput > lenExpected {
		fitness += (1.0 - (lenOutput / lenExpected))
	}

	if fitness > 0 {
		factor *= 10
		weight := 1.0 / float64(len(e.program)+1)
		fitness += (weight / factor)
	}

	e.fitness = fitness
}

func (e *Entry) String() string {
	return fmt.Sprintf("output = %#v fitness = %f runtime = %d generation = %d", string(e.output), e.fitness, e.runtime, e.generation)
}

type byFitness []Entry

func (e byFitness) Len() int      { return len(e) }
func (e byFitness) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e byFitness) Less(i, j int) bool {
	return e[i].fitness > e[j].fitness
}
