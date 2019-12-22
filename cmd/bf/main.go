package main

import (
	"bufio"
	"log"
	"os"

	"github.com/sanderhahn/go-bf"
)

func main() {
	for _, filename := range os.Args[1:] {
		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		i := bf.NewInterpreter(os.Stdout, os.Stdin)
		input := bufio.NewReader(file)
		if err := i.Interpret(input); err != nil {
			log.Fatal(err)
		}
	}
}
