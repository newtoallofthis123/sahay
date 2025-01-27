package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/newtoallofthis123/sahay/model"
	"github.com/newtoallofthis123/sahay/parser"
	"github.com/sbabiv/xml2map"
)

func main() {
	var filename string
	flag.StringVar(&filename, "file", "", "The file to read")
	flag.Parse()

	p, err := parser.NewParser(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(p.Tags)

	a, err := model.NewModelApi(nil)
	if err != nil {
		panic(err)
	}
	x, err := a.GetResponse(p.File.Seek[uint16(p.Tags[2].Line)])
	if err != nil {
		panic(err)
	}
	decoder := xml2map.NewDecoder(strings.NewReader(x))
	result, err := decoder.Decode()
	fmt.Println(result["Comment"])
}
