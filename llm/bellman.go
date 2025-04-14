package llm

import (
	"fmt"
	"strings"

	"github.com/modfin/bellman/models/gen"
	"github.com/modfin/bellman/prompt"
	"github.com/modfin/bellman/schema"
)

type BellmanProvider struct {
	generator *gen.Generator
}

func newBellmanProvider(generator *gen.Generator) *BellmanProvider {
	return &BellmanProvider{generator: generator}
}

func (p *BellmanProvider) GenerateTitle(commits []string) *string {
	if p.generator == nil {
		return nil
	}

	res, err := p.generator.
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

func (p *BellmanProvider) GenerateBody(commits []string, diff string, template string) *string {
	if p.generator == nil {
		return nil
	}

	var res *gen.Response
	var err error
	commitsJoined := strings.Join(commits, "\n")

	if template != "" {
		res, err = p.generator.
			Output(schema.From(BodyResponse{})).
			Prompt(prompt.AsUser(fmt.Sprintf(BODY_PROMPT_WITH_TEMPLATE, template, commitsJoined, diff)))
	} else {
		res, err = p.generator.
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
