package main

import (
	"bufio"
	"github.com/sanderhahn/bf"
	"log"
	"os"
)

func main() {
	for _, filename := range os.Args[1:] {
		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		i := bf.NewInterpreter(os.Stdin, os.Stdout)
		input := bufio.NewReader(file)
		if err := i.Interpret(input); err != nil {
			log.Fatal(err)
		}
	}
}
