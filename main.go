package main

import (
	"flag"
	"fmt"

	"github.com/newtoallofthis123/sahay/handler"
	"github.com/newtoallofthis123/sahay/model"
	"github.com/newtoallofthis123/sahay/parser"
)

func main() {
	var filename string
	var modelName string
	var prompt string
	flag.StringVar(&filename, "file", "", "The file to read")
	flag.StringVar(&modelName, "model", "mistral:7b", "The Model to be used (should be downloaded)")
	flag.StringVar(&prompt, "prompt", "", "The prompt to the used")
	flag.Parse()

	if filename == "" {
		fmt.Println("Please provide a file to read")
		return
	}

	p, err := parser.NewParser(filename)
	if err != nil {
		panic(err)
	}

	a, err := model.NewModelApi(&model.ModelOptions{Model: modelName, Prompt: prompt})
	if err != nil {
		panic(err)
	}

	h := handler.NewHandler(a)
	res, err := h.GetResponses(&p)
	if err != nil {
		panic(err)
	}

	err = h.WriteToFile(&p, res)
	if err != nil {
		panic(err)
	}

	fmt.Println("Modified file:", p.Filename)
}
