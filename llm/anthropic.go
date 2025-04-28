package llm

import (
	"context"
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	client *anthropic.Client
	modelName string
}

func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &AnthropicProvider{
		&client,
		"claude-3-5-haiku-latest",
	}
}

func (p *AnthropicProvider) GenerateTitleAndBody(commits []string, diff string, prTemplate string, ctx context.Context) (*string, *string) {
	prompt := getUserPrompt(commits, diff, prTemplate)

	schema := generateSchema[Response]()

	tool := anthropic.ToolParam{
		Name:        "generate_pr_description",
		Description: anthropic.String("Represents the description of a pull request with title and body"),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: schema.Properties,
		},
	}

	message, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     p.modelName,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: SYSTEM_PROMPT},
		},
		Messages: []anthropic.MessageParam{{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.TextBlockParam{Text: prompt},
			}},
		}},
		Tools: []anthropic.ToolUnionParam{
			{
				OfTool: &tool,
			},
		},
	})

	if err != nil {
		return nil, nil
	}

	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.ToolUseBlock:
			switch block.Name {
			case "generate_pr_description":
				var response Response

				err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &response)
				if err != nil {
					return nil, nil
				}

				return &response.Title, &response.Body
			}
		}

	}

	return nil, nil
}
