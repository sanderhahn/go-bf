package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sanderhahn/go-bf"
)

const iterations = 10000

func wrapAt(s string, at int) string {
	b := bytes.NewBufferString("")
	for pos := 0; pos < len(s); pos += at {
		if pos+at < len(s) {
			b.WriteString(s[pos : pos+at])
			b.WriteRune('\n')
		} else {
			b.WriteString(s[pos:])
		}
	}
	return b.String()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	expected, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	population := bf.NewPopulation()
	population.Expected = expected

	flag.IntVar(&population.MaxRuntime, "runtime", 10000, "max runtime for program")
	flag.IntVar(&population.MaxManipulate, "manipulate", 3, "max manipulation when copying")
	flag.Parse()

	for i := 1; i <= iterations; i++ {
		population.EvaluateAndMutate()
		fittest := population.Fittest()
		fmt.Printf("%d: %s\n", i, fittest)
		if code, ok := population.SuccessCode(); ok {
			fmt.Printf("%s\n", wrapAt(string(code), 80))
		}
	}
}
