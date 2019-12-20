package bf

import (
	"errors"
	"io"
)

const memorySize = 1024
const stackSize = 1024

var errInvalidNesting = errors.New("Invalid loop nesting")
var errMemory = errors.New("Memory below zero unsupported")

// Interpreter has the state for a single interpreter
type Interpreter struct {
	ptr    int
	memory []byte
	w      io.Writer
	r      io.Reader
	ip     int    // instruction pointer
	code   []byte // code memory
	stack  []int  // return stack
}

func memory() []byte {
	return make([]byte, memorySize)
}

// NewInterpreter constructor
func NewInterpreter(w io.Writer, r io.Reader) *Interpreter {
	return &Interpreter{
		ptr:    0,
		memory: memory(),
		w:      w,
		r:      r,
		ip:     0,
		code:   make([]byte, 0, memorySize),
		stack:  make([]int, 0, stackSize),
	}
}

func readByte(r io.Reader) (byte, error) {
	var buf = []byte{0}
	n, err := r.Read(buf)
	if n == 0 {
		return 0, io.EOF
	}
	return buf[0], err
}

func (i *Interpreter) skipLoop(r io.Reader) error {
	level := 0
	for {
		code, err := i.instr(r)
		if err != nil {
			return err
		}
		if code == '[' {
			level++
		} else if code == ']' {
			if level == 0 {
				break
			}
			if level > 0 {
				level--
			}
		}
	}
	return nil
}

func writeByte(w io.Writer, b byte) error {
	n, err := w.Write([]byte{b})
	if n != 1 || err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) increaseMemory() error {
	// increase memory on demand
	if i.ptr >= len(i.memory) {
		i.memory = append(i.memory, memory()...)
	}
	return nil
}

// instr reads and caches code instructions from the reader into code memory
func (i *Interpreter) instr(r io.Reader) (code byte, err error) {
	if i.ip < len(i.code) {
		// previously read code
		code = i.code[i.ip]
	} else {
		// read new instruction
		code, err = readByte(r)
		if err != nil {
			return
		}
		// store in code memory
		i.code = append(i.code, code)
	}
	i.ip++
	return
}

func (i *Interpreter) condition() bool {
	return i.memory[i.ptr] == 0
}

// push return position
func (i *Interpreter) push() {
	i.stack = append(i.stack, i.ip)
}

// jump to return position
func (i *Interpreter) jump() {
	top := len(i.stack) - 1
	i.ip = i.stack[top]
}

// pop return position from the stack
func (i *Interpreter) pop() {
	top := len(i.stack) - 1
	i.stack = i.stack[0:top]
}

// Interpret the instructions from the reader
func (i *Interpreter) Interpret(r io.Reader) error {
	for {
		code, err := i.instr(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch code {
		case '>':
			i.ptr++
			if err := i.increaseMemory(); err != nil {
				return err
			}
		case '<':
			i.ptr--
			if i.ptr < 0 {
				return errMemory
			}
		case '+':
			i.memory[i.ptr]++
		case '-':
			i.memory[i.ptr]--
		case '.':
			b := i.memory[i.ptr]
			err := writeByte(i.w, b)
			if err != nil {
				return err
			}
		case ',':
			b, err := readByte(i.r)
			if err == io.EOF {
				// dbf2c.bf expects zero as EOF
				b = 0
			} else if err != nil {
				return err
			}
			i.memory[i.ptr] = b
		case '[':
			if i.condition() {
				if err := i.skipLoop(r); err != nil {
					return err
				}
			} else {
				i.push()
			}
		case ']':
			if len(i.stack) < 1 {
				return errInvalidNesting
			}
			if i.condition() {
				i.pop()
			} else {
				i.jump()
			}
		default:
			// ignore comments
		}
	}
	return nil
}
