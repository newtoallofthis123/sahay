package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/newtoallofthis123/sahay/file"
)

type Parser struct {
	Filename string
	File     file.File
	Tags     []Tag
}

type Tag struct {
	Type      string `json:"_type,omitempty"`
	Name      string `json:"name,omitempty"`
	FilePath  string `json:"path,omitempty"`
	Language  string `json:"language,omitempty"`
	Line      int    `json:"line,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	TagKind   string `json:"kind,omitempty"`
	Scope     string `json:"scope,omitempty"`
	ScopeKind string `json:"scopeKind,omitempty"`
}

func NewParser(filename string) (Parser, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return Parser{}, err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return Parser{}, fmt.Errorf("file does not exist: %s", filename)
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return Parser{}, err
	}

	fileRepr, err := file.NewFile(f)
	if err != nil {
		return Parser{}, err
	}

	ext := strings.Replace(filepath.Ext(filename), ".", "", 1)

	rawLines, err := RunCtags(filename, &CTagsOptions{LangExtension: ext})
	if err != nil {
		return Parser{}, err
	}

	tags := make([]Tag, 0)

	for _, raw := range rawLines {
		if raw == "" || strings.Contains(raw, "Warning") {
			continue
		}
		tag, err := toTag(raw)
		if err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}

	return Parser{
		Tags:     tags,
		Filename: filename,
		File:     fileRepr,
	}, nil
}

func toTag(raw string) (Tag, error) {
	var tag Tag
	err := json.Unmarshal([]byte(raw), &tag)
	return tag, err
}
