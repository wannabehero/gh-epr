package llm

import (
	"context"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client    *genai.Client
	modelName string
}

func NewGeminiProvider(apiKey string, ctx context.Context) *GeminiProvider {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil
	}

	return &GeminiProvider{
		client,
		"gemini-2.5-flash-preview-04-17",
	}
}

func (p *GeminiProvider) newModelWithSchema(schema *genai.Schema) *genai.GenerativeModel {
	model := p.client.GenerativeModel(p.modelName)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = schema
	return model
}

func (p *GeminiProvider) GenerateTitleAndBody(commits []string, diff string, prTemplate string, ctx context.Context) (*string, *string) {
	model := p.newModelWithSchema(&genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"body":  {Type: genai.TypeString},
			"title": {Type: genai.TypeString},
		},
		Required: []string{"body", "title"},
	})
	model.SystemInstruction = genai.NewUserContent(genai.Text(SYSTEM_PROMPT))

	prompt := getUserPrompt(commits, diff, prTemplate)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil || len(resp.Candidates) == 0 {
		return nil, nil
	}

	var text string
	if resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				text = string(txt)
			}
		}
	}

	var response Response
	if err := json.Unmarshal([]byte(text), &response); err != nil {
		return nil, nil
	}

	var body, title *string
	if response.Body != "" {
		body = &response.Body
	}
	if response.Title != "" {
		title = &response.Title
	}

	return title, body
}
