package handler

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/newtoallofthis123/sahay/model"
	"github.com/newtoallofthis123/sahay/parser"
	"github.com/sbabiv/xml2map"
)

type Handler struct {
	Api *model.ModelApi
}

func NewHandler(api *model.ModelApi) Handler {
	return Handler{
		Api: api,
	}
}

// parseRes parses the given XML string and returns a result string or an error.
func (h *Handler) parseRes(xml string) (string, error) {
	decoder := xml2map.NewDecoder(strings.NewReader(xml))
	result, err := decoder.Decode()
	if err != nil {
		return "", err
	}

	res, ok := result["Comment"].(string)
	if !ok {
		return "", fmt.Errorf("Result is not a string")
	}

	return res, nil
}

// TODO: Make it run in parallel
func (h *Handler) GetResponses(p *parser.Parser) (map[uint16]string, error) {

	res := make(map[uint16]string, 0)

	for _, tag := range p.Tags {
		raw := p.File.Seek[uint16(tag.Line)]
		fmt.Print("Adding comment for function:", raw)
		xml, err := h.Api.GetResponse(raw)
		if err != nil {
			return nil, err
		}

		comment, err := h.parseRes(xml)
		if err != nil {
			return nil, err
		}

		res[uint16(tag.Line)] = comment
	}

	return res, nil
}

func (h *Handler) WriteToFile(p *parser.Parser, comments map[uint16]string) error {
	file, err := os.OpenFile(p.File.GetName(), os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	keys := make([]uint16, 0, len(p.File.Seek))
	for k := range p.File.Seek {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, num := range keys {
		line := p.File.Seek[num]
		comment, ok := comments[num]
		if ok {
			t := line
			line = ""
			// FIXME: It has to be based on the language
			if !strings.Contains(comment, "//") {
				line += "//"
			}
			line += comment
			line += "\n" + t
		}
		io.WriteString(file, line)
	}

	return err
}
