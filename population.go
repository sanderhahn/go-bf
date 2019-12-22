package bf

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strings"
)

func (e *Entry) exec(input string, maxRuntime int) error {
	output := &strings.Builder{}
	i := NewInterpreter(output, bytes.NewReader([]byte(input)))
	runtime, err := i.InterpretExtended(bytes.NewReader(e.program), false, maxRuntime)
	if err == nil {
		e.runtime = maxRuntime - runtime
		e.output = output.String()
	} else {
		e.runtime = 0
		e.output = ""
	}
	return e.err
}

func (e *Entry) calculateFitness(expected string) {
	fitness := 0.0
	for i, ch := range []byte(expected) {
		if i < len(e.output) {
			diff := math.Abs(float64(ch)-float64(e.output[i])) / 255.0
			fitness += (1.0 - diff)
			fitness *= 3.0
		}
	}

	lenOutput := float64(len(e.output))
	lenExpected := float64(len(expected))
	if lenOutput > lenExpected {
		fitness *= (1.0 - (lenOutput / lenExpected))
	}

	if e.success {
		// optimize for program length once successful
		weight := 1.0 / float64(len(e.program))
		fitness += weight
	}

	e.fitness = fitness
}

const keepSize = 32
const manipulationSize = 32
const populationSize = keepSize * manipulationSize

// Population maintains the pool of programs in the form of entries
type Population struct {
	entries    []Entry
	Expected   string
	MaxRuntime int
}

// Entry maintains information of a program
type Entry struct {
	program    Program
	output     string
	runtime    int
	err        error
	fitness    float64
	success    bool
	generation int
}

func (e *Entry) String() string {
	return fmt.Sprintf("output = %#v fitness = %f runtime = %d generation = %d", e.output, e.fitness, e.runtime, e.generation)
}

type byFitness []Entry

func (e byFitness) Len() int      { return len(e) }
func (e byFitness) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e byFitness) Less(i, j int) bool {
	return e[i].fitness > e[j].fitness
}

// NewPopulation constructor
func NewPopulation() *Population {
	entries := make([]Entry, populationSize)
	for i := range entries {
		entries[i].program = NewRandomProgram(1)
	}
	return &Population{
		entries: entries,
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
		entry.success = entry.err == nil && entry.output == p.Expected
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
				entry.program = Mutate(entry.program, p.entries[0:keepSize])
				entry.generation = keepEntry.generation
			}
		}
	}
}
