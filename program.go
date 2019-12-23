package bf

import (
	"bytes"
	"math/rand"
)

// Instr represents a instruction
type Instr byte

const instr = `><+-.[]`

// RandomInstr returns a random instruction
func RandomInstr() byte {
	return instr[rand.Intn(len(instr))]
}

// Program represents a program
type Program []byte

// NewProgram constructs empty program
func NewProgram() Program {
	return make([]byte, 0, 128)
}

// NewProgramClone clones a program
func NewProgramClone(source Program) Program {
	return append(NewProgram(), source...)
}

// NewRandomProgram makes a random program
func NewRandomProgram(length int) Program {
	p := NewProgram()
	for i := 0; i < length; i++ {
		p = append(p, RandomInstr())
	}
	return p
}

// Normalize the code to be syntactically valid with respect to loops
func Normalize(program Program) Program {
	buf := bytes.NewBuffer(Program{})
	level := 0
	// processed and written backwards
	for i := range program {
		r := len(program) - i - 1
		b := program[r]
		if b == ']' {
			level++
		} else if b == '[' {
			if level <= 0 {
				continue
			}
			level--
		}
		buf.WriteByte(b)
	}
	reversed := buf.Bytes()
	buf = bytes.NewBuffer(Program{})
	level = 0
	// processed in reverse written forwards
	for i := range reversed {
		r := len(reversed) - i - 1
		b := reversed[r]
		if b == '[' {
			level++
		} else if b == ']' {
			if level <= 0 {
				continue
			}
			level--
		}
		buf.WriteByte(b)
	}
	return buf.Bytes()
}

func insertAt(a Program, pos int, b Program) Program {
	return append(a[:pos], append(b, a[pos:]...)...)
}

func removeAt(p Program, pos, len int) Program {
	return append(p[0:pos], p[pos+len:]...)
}

func replaceAt(a Program, pos int, len int, b Program) Program {
	c := removeAt(a, pos, len)
	return insertAt(c, pos, b)
}

var compounds = []Program{
	Program(`[-]`),                     // set zero
	Program(`>++++[<++++>-]<`),         // add 0x10
	Program(`>++++[<++++++++>-]<`),     // add 0x20
	Program(`>++++++++[<++++++++>-]<`), // add 0x40
}

// Loop a program
func (p Program) Loop() Program {
	return append(append(Program(`[`), p...), ']')
}

// Comment out a program
func (p Program) Comment() Program {
	return append(append(Program(`[-][`), p...), ']')
}

// Mutate a program randomly a number of times
func Mutate(code Program, times int, sources []Entry) Program {
	for i := 0; i < times; i++ {
		code = mutate(code, sources)
	}
	return code
}

func mutate(code Program, sources []Entry) Program {
	if len(code) == 0 {
		return NewRandomProgram(1)
	}

	pos := rand.Intn(len(code))
	length := rand.Intn(len(code) - pos)

	switch rand.Intn(8) {
	case 0:
		return insertAt(code, pos, NewRandomProgram(1))
	case 1:
		return removeAt(code, pos, 1)
	case 2:
		p := NewProgramClone(code)
		p[pos] = RandomInstr()
		return p
	case 3:
		apos := rand.Intn(len(code))
		len := rand.Intn(len(code) - apos)
		return insertAt(code, pos, code[apos:apos+len])
	case 4:
		// cross breed
		pick := rand.Intn(len(sources))
		return insertAt(code, pos, sources[pick].program)
	case 5:
		return removeAt(code, pos, length)
	case 6:
		// insert some loops for variety
		c := rand.Intn(len(compounds))
		return insertAt(code, pos, compounds[c])
	case 7:
		without := removeAt(code, pos, length)
		return insertAt(without, pos, code[pos:pos+length].Comment())
	}
	panic("unreachable")
}
