package llm

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenaiProvider struct {
	client *openai.Client
}

func generateSchema[T any]() *jsonschema.Schema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func NewOpenaiProvider(apiKey string) *OpenaiProvider {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &OpenaiProvider{
		&client,
	}
}

func (p *OpenaiProvider) GenerateTitleAndBody(commits []string, diff string, prTemplate string, ctx context.Context) (*string, *string) {
	schema := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "pr_description",
		Description: openai.String("Description of a PR"),
		Schema:      generateSchema[Response](),
		Strict:      openai.Bool(true),
	}

	prompt := getUserPrompt(commits, diff, prTemplate)

	chat, _ := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage(SYSTEM_PROMPT),
			openai.UserMessage(prompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schema,
			},
		},
		Model: openai.ChatModelGPT4o,
	})

	var response Response
	_ = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &response)
	return &response.Title, &response.Body
}
