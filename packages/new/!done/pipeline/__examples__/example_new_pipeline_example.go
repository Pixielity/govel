package main

import (
	"fmt"
	"govel/new/pipeline"
)

// ExampleNewPipeline demonstrates creating a new Pipeline instance.
func ExampleNewPipeline() {
	p := pipeline.NewPipeline(nil)
	result, _ := p.Send("hello").ThenReturn()
	fmt.Println(result)
	// Output: hello
}

func main() {
	ExampleNewPipeline()
}

