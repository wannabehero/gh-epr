package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const SYSTEM_PROMPT = `You are a tool that summarises git commit messages
and git diffs to generate a descriptive Pull Request title and body.

Title describes the changes in no more than 10 words.
Start with a verb like "create", "fix", "update", etc.
Put a random emoji that describes the context in front of the title.
For instance: '👾 Update the build CI workflow'

Your responses are concise and to the point`

const BODY_TEMPLATE_PROMPT = `IMPORTANT FOR THE BODY: Format your response according to the provided PR template structure.
Fill in the sections appropriately while maintaining the template format.

PR Template:
%s`

const USER_PROMPT = `
Using the following commit messages and a diff as a context
generate a descriptive consise Pull Request title and body that summarizes the changes.

For the title include an emoji that describes the context to the start.

For the body include a brief summary section with 1-3 bullet points describing the key changes.
Format it in Markdown with proper headings.

%s

Commit messages:
%s

Diff:
%s
`

func getUserPrompt(commits []string, diff string, prTemplate string) string {
	var prTermplatePrompt string
	if prTemplate != "" {
		prTermplatePrompt = fmt.Sprintf(BODY_TEMPLATE_PROMPT, prTemplate)
	} else {
		prTermplatePrompt = ""
	}

	prompt := fmt.Sprintf(USER_PROMPT, prTermplatePrompt, strings.Join(commits, "\n"), diff)
	return prompt
}

func getSystemPrompt() string {
	basePrompt := SYSTEM_PROMPT

	if override := viper.GetString("system_prompt_override"); override != "" {
		basePrompt = override
	}

	if instructions := viper.GetString("additional_instructions"); instructions != "" {
		return basePrompt + "\n\n## Additional User-Customized Instructions\n" + instructions
	}

	return basePrompt
}

type Response struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func getProvider(ctx context.Context) LLMProvider {
	providerCreators := map[string]func(context.Context) LLMProvider{
		"gemini": func(ctx context.Context) LLMProvider {
			key := os.Getenv("GEMINI_API_KEY")
			if key == "" {
				return nil
			}
			return NewGeminiProvider(key, ctx)
		},
		"openai": func(ctx context.Context) LLMProvider {
			key := os.Getenv("OPENAI_API_KEY")
			if key == "" {
				return nil
			}
			return NewOpenaiProvider(key)
		},
		"anthropic": func(ctx context.Context) LLMProvider {
			key := os.Getenv("ANTHROPIC_API_KEY")
			if key == "" {
				return nil
			}
			return NewAnthropicProvider(key)
		},
	}

	configuredProvider := viper.GetString("provider")

	providersToTry := []string{"gemini", "openai", "anthropic"}
	if configuredProvider != "" {
		if _, exists := providerCreators[configuredProvider]; !exists {
			fmt.Fprintf(os.Stderr, "Error: Unknown provider '%s' specified in config\n", configuredProvider)
			return nil
		}
		providersToTry = []string{configuredProvider}
	}

	for _, providerName := range providersToTry {
		createFn := providerCreators[providerName]
		if provider := createFn(ctx); provider != nil {
			return provider
		} else if configuredProvider != "" {
			fmt.Fprintf(os.Stderr, "Error: Provider '%s' specified in config but required API key is not set\n", providerName)
			return nil
		}
	}

	if configuredProvider == "" {
		fmt.Fprintf(os.Stderr, "Error: No API keys found for any supported provider\n")
	}
	return nil
}

func GenerateTitleAndBody(commits []string, diff string, template string, ctx context.Context) (*string, *string) {
	provider := getProvider(ctx)
	if provider == nil {
		return nil, nil
	}
	return provider.GenerateTitleAndBody(commits, diff, template, ctx)
}
