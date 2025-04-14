package llm

import (
	"os"

	"github.com/modfin/bellman/services/anthropic"
	"github.com/modfin/bellman/services/openai"
)

const SYSTEM_PROMPT = `You a git commit messages
summariser that produce playful Pull Reuqest titles
prefixed with an emoji.
Your responses are consise`

const COMMITS_PROMPT = `
Using the following commit messages as a context
generate a descriptive consise Pull Request title
that would describe the changes in no more than 10 words.
Start with a verb like "create", "fixe", "update", etc.
Put a random emoji that describes the context in front of the title.
For instance: '👾 Update the build CI workflow'

Commit messages:
%s
`

const BODY_PROMPT = `
Using the following commit messages and diff as context,
generate a descriptive Pull Request body that summarizes the changes.
Include a brief summary section with 1-3 bullet points describing the key changes.
Be very concise and to the point.

Format it in Markdown with proper headings.

Commit messages:
%s

Diff:
%s
`

const BODY_PROMPT_WITH_TEMPLATE = `
Using the following commit messages and diff as context,
generate a descriptive Pull Request body that summarizes the changes.
Be very concise and to the point.

IMPORTANT: Format your response according to the provided PR template structure.
Fill in the sections appropriately while maintaining the template format.

PR Template:
%s

Commit messages:
%s

Diff:
%s
`

type Response struct {
	Title string `json:"title"`
}

type BodyResponse struct {
	Body string `json:"body"`
}

func getProvider() LLMProvider {
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		provider := NewGeminiProvider(key)
		if provider != nil {
			return provider
		}
	}

	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		generator := openai.New(key).Generator().Model(openai.GenModel_gpt4o)
		return newBellmanProvider(generator)
	}

	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		generator := anthropic.New(key).Generator().Model(anthropic.GenModel_3_5_sonnet_latest)
		return newBellmanProvider(generator)
	}

	return nil
}

func GenerateTitle(commits []string) *string {
	provider := getProvider()
	if provider == nil {
		return nil
	}
	return provider.GenerateTitle(commits)
}

func GenerateBody(commits []string, diff string, template string) *string {
	provider := getProvider()
	if provider == nil {
		return nil
	}
	return provider.GenerateBody(commits, diff, template)
}
