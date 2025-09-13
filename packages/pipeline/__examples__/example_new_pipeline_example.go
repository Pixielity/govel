package main

import (
	"fmt"
	"govel/packages/new/pipeline/src"
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

