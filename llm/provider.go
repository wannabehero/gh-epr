package llm

type LLMProvider interface {
	GenerateTitleAndBody(commits []string, diff string, template string) (*string, *string)
}
