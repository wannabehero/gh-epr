package llm

import "context"

type LLMProvider interface {
	GenerateTitleAndBody(commits []string, diff string, template string, ctx context.Context) (*string, *string)
}
