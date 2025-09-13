package main

import (
	"fmt"
	"govel/packages/new/pipeline/src"
)

// ExamplePipeline_Through demonstrates using the Through method with function and object pipes.
func ExamplePipeline_Through() {
	p := pipeline.NewPipeline(nil)

	// Function pipe
	addSuffix := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
		if s, ok := passable.(string); ok {
			passable = s + "-fn"
		}
		return next(passable)
	}

	// Execute
	out, _ := p.
		Send("in").
		Through([]interface{}{addSuffix}).
		ThenReturn()

	fmt.Println(out)
	// Output: in-fn
}

func main() {
	ExamplePipeline_Through()
}

