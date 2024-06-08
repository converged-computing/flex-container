package main

import (
	"flag"
	"fmt"

	"github.com/converged-computing/flex-container/src/graph"
)

func main() {
	fmt.Println("This is the flex kubernetes prototype")
	matchPolicy := flag.String("policy", "first", "Match policy")
	nodelist := flag.String("nodelist", "graph.json", "Read nodelist from this file")

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"examples/jobspec-1.yaml"} //, "examples/jobspec-2.yaml"}
	}
	fmt.Println(args)

	// Create an ice cream graph, and match the spec to it.
	g := graph.NewClusterGraph(*matchPolicy)
	g.CreateGraph(*nodelist)
	fmt.Println(g)

	for _, jobspec := range args {
		err := g.Match(jobspec)
		if err != nil {
			fmt.Printf("There was a problem with your request: %s\n", err)
			return
		}
	}

}
