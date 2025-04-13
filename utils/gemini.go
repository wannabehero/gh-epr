package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Response struct {
	Title string `json:"title"`
}

type BodyResponse struct {
	Body string `json:"body"`
}

type GeminiProvider struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

const jsonInstructionTemplate = `Important: You must respond with ONLY valid JSON in the format: %s
Do not include any other text, markdown formatting, or explanations in your response.`

func NewGeminiProvider(apiKey string) *GeminiProvider {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil
	}
	model := client.GenerativeModel("gemini-2.0-flash")
	
	return &GeminiProvider{
		client: client,
		model:  model,
		ctx:    ctx,
	}
}

func (p *GeminiProvider) generateJSON(prompt string, jsonFormat string) (string, error) {
	if p.model == nil {
		return "", fmt.Errorf("model not initialized")
	}
	
	fullPrompt := fmt.Sprintf("%s\n%s", prompt, fmt.Sprintf(jsonInstructionTemplate, jsonFormat))
	p.model.SetResponseMIMEType("application/json")
	
	resp, err := p.model.GenerateContent(p.ctx, genai.Text(fullPrompt))
	if err != nil || len(resp.Candidates) == 0 {
		return "", err
	}
	
	var text string
	if resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			text += fmt.Sprintf("%v", part)
		}
	}
	
	return text, nil
}

func (p *GeminiProvider) GenerateTitle(commits []string) *string {
	prompt := fmt.Sprintf("%s\n%s", SYSTEM_PROMPT, fmt.Sprintf(COMMITS_PROMPT, strings.Join(commits, "\n")))
	
	text, err := p.generateJSON(prompt, `{"title": "your title here"}`)
	if err != nil {
		return nil
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
	var prompt string
	commitsJoined := strings.Join(commits, "\n")
	if template != "" {
		prompt = fmt.Sprintf(BODY_PROMPT_WITH_TEMPLATE, template, commitsJoined, diff)
	} else {
		prompt = fmt.Sprintf(BODY_PROMPT, commitsJoined, diff)
	}
	
	text, err := p.generateJSON(prompt, `{"body": "your body markdown here"}`)
	if err != nil {
		return nil
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
