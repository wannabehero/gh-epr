package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/modfin/bellman/models/gen"
	"github.com/modfin/bellman/prompt"
	"github.com/modfin/bellman/schema"
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
Be concise and to the point.

Format it in Markdown with proper headings.

Commit messages:
%s

Diff:
%s
`

const BODY_PROMPT_WITH_TEMPLATE = `
Using the following commit messages and diff as context,
generate a descriptive Pull Request body that summarizes the changes.

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

func getAvailableGenerator() *gen.Generator {
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return openai.New(key).Generator().Model(openai.GenModel_gpt4o)
	}

	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return anthropic.New(key).Generator().Model(anthropic.GenModel_3_5_sonnet_latest)
	}

	return nil
}

func GenerateTitle(commits []string) *string {
	generator := getAvailableGenerator()

	if generator == nil {
		return nil
	}

	res, err := generator.
		Output(schema.From(Response{})).
		Prompt(prompt.AsUser(fmt.Sprintf(COMMITS_PROMPT, strings.Join(commits, "\n"))))

	if err != nil {
		return nil
	}

	var response Response
	if err = res.Unmarshal(&response); err != nil {
		return nil
	}

	return &response.Title
}

func GenerateBody(commits []string, diff string, template string) *string {
	generator := getAvailableGenerator()

	if generator == nil {
		return nil
	}

	var res *gen.Response
	var err error
	commitsJoined := strings.Join(commits, "\n")

	if template != "" {
		res, err = generator.
			Output(schema.From(BodyResponse{})).
			Prompt(prompt.AsUser(fmt.Sprintf(BODY_PROMPT_WITH_TEMPLATE, template, commitsJoined, diff)))
	} else {
		res, err = generator.
			Output(schema.From(BodyResponse{})).
			Prompt(prompt.AsUser(fmt.Sprintf(BODY_PROMPT, commitsJoined, diff)))
	}

	if err != nil {
		return nil
	}

	var response BodyResponse
	if err = res.Unmarshal(&response); err != nil {
		return nil
	}

	return &response.Body
}
