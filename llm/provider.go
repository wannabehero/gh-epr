package llm

type LLMProvider interface {
	GenerateTitle(commits []string) *string
	GenerateBody(commits []string, diff string, template string) *string
}
