package openai

import (
	"context"

	"github.com/pkg/errors"
	openaigo "github.com/sashabaranov/go-openai"
	"github.com/vipulvpatil/airetreat-go/internal/utilities/logger"
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
	logger.LogMessageln(prompt)
	openAiGoClient := openaigo.NewClient(c.apiKey)
	ctx := context.Background()

	req := openaigo.CompletionRequest{
		Model:     openaigo.GPT3TextDavinci003,
		MaxTokens: 50,
		Prompt:    prompt,
	}
	resp, err := openAiGoClient.CreateCompletion(ctx, req)
	if err != nil {
		logger.LogMessagef("Completion error: %v\n", err)
		return "", errors.Wrap(err, "Open Ai error")
	}
	return resp.Choices[0].Text, nil
}
