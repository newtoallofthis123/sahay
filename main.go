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
	flag.StringVar(&filename, "file", "", "The file to read")
	flag.Parse()

	p, err := parser.NewParser(filename)
	if err != nil {
		panic(err)
	}

	a, err := model.NewModelApi(nil)
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
