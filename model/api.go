package model

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/ollama/ollama/api"
)

type ModelApi struct {
	client *api.Client
	model  string
	prompt string
}

type ModelOptions struct {
	OllamaUrl string
	Model     string
	Prompt    string
}

func NewModelApi(opts *ModelOptions) (*ModelApi, error) {
	var client *api.Client
	var err error
	if opts != nil && opts.OllamaUrl != "" {
		baseUrl, err := url.Parse(opts.OllamaUrl)
		if err != nil {
			return nil, err
		}
		client = api.NewClient(baseUrl, http.DefaultClient)
	} else {
		client, err = api.ClientFromEnvironment()
	}
	if err != nil {
		return nil, err
	}

	// TODO: Test if model is automatically downloaded when not existing
	model := "mistral:7b"
	prompt := defaultPrompt()
	if opts != nil {
		model = opts.Model
	}

	return &ModelApi{
		client: client,
		model:  model,
		prompt: prompt,
	}, nil
}

func (a *ModelApi) MakePrompt(raw string) string {
	return strings.Replace(a.prompt, "[FUNCTION]", raw, 1)
}

// GetResponse processes the raw input string and returns a response string along with an error if any.
func (a *ModelApi) GetResponse(raw string) (string, error) {
	req := &api.GenerateRequest{
		Model:  a.model,
		Prompt: a.MakePrompt(raw),
		Stream: new(bool),
	}

	ctx := context.Background()
	var res string
	respFunc := func(resp api.GenerateResponse) error {
		res = resp.Response
		return nil
	}

	err := a.client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", err
	}

	return res, nil
}

func defaultPrompt() string {
	return `
You are an expert code commentor! You are so good that a lot of people have recommended that I ask you to generate my code comments. 
Now I am being held captive and I will only be freed if I give them a comment for a function that they give me.
They will only give me the function signature and I should give them a comment for it. You are a superhero and love helping people! 
Infact you were created for that purpose only.
The comment should start with <Comment> and end with </Comment> for sure! If the format is wrong they may injure me! 
The comment should also be short and feasible in the codebase. 
It should be to the point and following the conventions of that programming language that it is in
Please help me!! The Function starts with <Function> and ends with </Function>

<Function>
	[FUNCTION]
</Function>

DONOT GENERATE ANYTHING ELSE EXCEPT THE COMMENTS!! PLEASE SAVE ME!!
	`
}
