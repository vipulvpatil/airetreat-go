package openai

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	openaigo "github.com/sashabaranov/go-openai"
)

type client struct {
	apiKey string
}

type Client interface {
	CallCompletionApi(prompt string) (string, error)
}

type OpenAiClientOptions struct {
	ApiKey string
}

func NewClient(opts OpenAiClientOptions) Client {
	return &client{
		apiKey: opts.ApiKey,
	}
}

func (c *client) CallCompletionApi(prompt string) (string, error) {
	fmt.Println(prompt)
	openAiGoClient := openaigo.NewClient(c.apiKey)
	ctx := context.Background()

	req := openaigo.CompletionRequest{
		Model:     openaigo.GPT3TextDavinci003,
		MaxTokens: 5,
		Prompt:    prompt,
	}
	resp, err := openAiGoClient.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return "", errors.Wrap(err, "Open Ai error")
	}
	fmt.Println(resp.Choices[0])
	return resp.Choices[0].Text, nil
}
