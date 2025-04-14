package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil
	}
	model := client.GenerativeModel("gemini-2.0-flash")

	model.ResponseMIMEType = "application/json"

	return &GeminiProvider{
		client: client,
		model:  model,
		ctx:    ctx,
	}
}

func (p *GeminiProvider) GenerateTitle(commits []string) *string {
	if p.model == nil {
		return nil
	}

	p.model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {Type: genai.TypeString},
		},
		Required: []string{"title"},
	}

	prompt := fmt.Sprintf("%s\n%s", SYSTEM_PROMPT, fmt.Sprintf(COMMITS_PROMPT, strings.Join(commits, "\n")))

	fullPrompt := fmt.Sprintf("%s\nRespond with a JSON object containing only a 'title' field.", prompt)

	resp, err := p.model.GenerateContent(p.ctx, genai.Text(fullPrompt))
	if err != nil || len(resp.Candidates) == 0 {
		return nil
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
		return nil
	}

	if response.Title == "" {
		return nil
	}

	return &response.Title
}

func (p *GeminiProvider) GenerateBody(commits []string, diff string, template string) *string {
	if p.model == nil {
		return nil
	}

	p.model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"body": {Type: genai.TypeString},
		},
		Required: []string{"body"},
	}

	var prompt string
	commitsJoined := strings.Join(commits, "\n")
	if template != "" {
		prompt = fmt.Sprintf(BODY_PROMPT_WITH_TEMPLATE, template, commitsJoined, diff)
	} else {
		prompt = fmt.Sprintf(BODY_PROMPT, commitsJoined, diff)
	}

	fullPrompt := fmt.Sprintf("%s\nRespond with a JSON object containing only a 'body' field.", prompt)

	resp, err := p.model.GenerateContent(p.ctx, genai.Text(fullPrompt))
	if err != nil || len(resp.Candidates) == 0 {
		return nil
	}

	var text string
	if resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				text = string(txt)
			}
		}
	}

	var response BodyResponse
	if err := json.Unmarshal([]byte(text), &response); err != nil {
		return nil
	}

	if response.Body == "" {
		return nil
	}

	return &response.Body
}
