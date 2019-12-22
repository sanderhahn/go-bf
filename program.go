package bf

import (
	"bytes"
	"math/rand"
)

// Program represents a program
type Program []byte

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

var instr = `><+-.[]`

func randomInstr() byte {
	return instr[rand.Intn(len(instr))]
}

func randomCode(length int) Program {
	buf := make([]byte, 0, 128)
	for i := 0; i < length; i++ {
		buf = append(buf, randomInstr())
	}
	return buf
}

func insertAt(code Program, pos int, program Program) Program {
	return append(code[:pos], append(program, code[pos:]...)...)
}

func removeAt(code Program, pos, len int) Program {
	return append(code[0:pos], code[pos+len:]...)
}

var componds = []Program{
	Program(`[-]`),                     // set zero
	Program(`>++++[<++++>-]<`),         // add 0x10
	Program(`>++++[<++++++++>-]<`),     // add 0x20
	Program(`>++++++++[<++++++++>-]<`), // add 0x40
}

func comment(code Program) Program {
	return append(append(Program(`[-][`), code...), ']')
}

func mutate(code Program, keep []Entry) Program {
	if len(code) == 0 {
		return randomCode(1)
	}
	pos := rand.Intn(len(code))

	switch rand.Intn(8) {
	case 0:
		return insertAt(code, pos, Program{randomInstr()})
	case 1:
		if len(code) > 1 {
			return removeAt(code, pos, 1)
		}
		return randomCode(1)
	case 2:
		code[pos] = randomInstr()
		return code
	case 3:
		apos := rand.Intn(len(code))
		len := rand.Intn(len(code) - apos)
		return insertAt(code, pos, code[apos:apos+len])
	case 4:
		// cross breed
		pick := rand.Intn(len(keep))
		return insertAt(code, pos, keep[pick].program)
	case 5:
		len := rand.Intn(len(code) - pos)
		return removeAt(code, pos, len)
	case 6:
		// insert some loops for variety
		c := rand.Intn(len(componds))
		return insertAt(code, pos, componds[c])
	case 7:
		// comment out selection instead of delete
		len := rand.Intn(len(code) - pos)
		commented := comment(code[pos : pos+len])
		without := removeAt(code, pos, len)
		return insertAt(without, pos, commented)
	}
	return code
}
