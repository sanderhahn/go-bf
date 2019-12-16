package bf

import (
	"bytes"
	"errors"
	"io"
)

const memorySize = 2048
const bufSize = 64

var errInvalidNesting = errors.New("Invalid loop nesting")
var errMemoryError = errors.New("Invalid memory access")

// Interpreter has the state for a single interpreter
type Interpreter struct {
	ptr    int
	memory []byte
	w      io.Writer
	r      io.Reader
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
	}
}

func readByte(r io.Reader) (byte, error) {
	var buf = []byte{0}
	n, err := r.Read(buf)
	if n == 0 {
		return 0, io.EOF
	}
	// log.Printf("byte = %c", buf[0])
	return buf[0], err
}

func readLoop(r io.Reader) ([]byte, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, bufSize))
	level := 0
	for {
		code, err := readByte(r)
		if err != nil {
			return nil, err
		}
		if code == '[' {
			level++
		} else if code == ']' {
			if level > 0 {
				level--
			} else {
				break
			}
		}
		buf.WriteByte(code)
	}
	return buf.Bytes(), nil
}

func writeByte(w io.Writer, b byte) error {
	n, err := w.Write([]byte{b})
	if n != 1 || err != nil {
		return err
	}
	// log.Printf("write %c", b)
	return nil
}

func (i *Interpreter) checkMemory() {
	// increase memory on demand
	if i.ptr >= len(i.memory) {
		i.memory = append(i.memory, memory()...)
		// log.Printf("memory increased to %d\n", len(i.memory))
	}
}

// Interpret the instructions from the reader
func (i *Interpreter) Interpret(r io.Reader) error {
	for {
		code, err := readByte(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch code {
		case '>':
			i.ptr++
		case '<':
			i.ptr--
		case '+':
			i.checkMemory()
			i.memory[i.ptr]++
		case '-':
			if i.ptr < 0 {
				return errMemoryError
			}
			i.memory[i.ptr]--
		case '.':
			i.checkMemory()
			b := i.memory[i.ptr]
			err := writeByte(i.w, b)
			if err != nil {
				return err
			}
		case ',':
			b, err := readByte(i.r)
			if err != nil {
				return err
			}
			i.memory[i.ptr] = b
		case '[':
			loop, err := readLoop(r)
			if err != nil {
				return err
			}
			for i.memory[i.ptr] != 0 {
				err := i.Interpret(bytes.NewReader(loop))
				if err != nil {
					return err
				}
			}
		case ']':
			return errInvalidNesting
		default:
			// ignore comments
		}
	}
	return nil
}
