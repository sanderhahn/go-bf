package main

import (
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

func main() {
	rand.Seed(time.Now().UnixNano())

	expected, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	population := bf.NewPopulation()
	population.Expected = string(expected)

	flag.IntVar(&population.MaxRuntime, "runtime", 10000, "max runtime for program")
	flag.IntVar(&population.MaxManipulation, "manipulation", 1, "max manipulation when copying")
	flag.Parse()

	for i := 1; i <= iterations; i++ {
		population.EvaluateAndMutate()
		fittest := population.Fittest()
		fmt.Printf("%d: %s\n", i, fittest)
		if code, ok := population.SuccessCode(); ok {
			fmt.Printf("%s\n", code)
		}
	}
}
