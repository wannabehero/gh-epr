# gh-aipr

Github CLI extension to generate pull requests titles and descriptions
based on your current branch commits and the diff.

## Installation

Make sure the CLI is installed: https://cli.github.com/

```
gh extension install wannabehero/gh-aipr
```

## Configuration

To make the best use of the tool set one of the following:
- `OPENAI_API_KEY`
- `ANTHROPIC_API_KEY`
- `GEMINI_API_KEY`

in your environment variables so it can use the LLM
to generate relevant title automatically
based on your current branch commits.

The tool uses some sensible defaults for the models
but you can override them in the config file.

See more in [docs/config.md](docs/config.md).

## Usage

```
gh aipr <other args that gh pr create supports>
```

This will follow you through interactive create process
and your final PR title and body will be automatically generated and set.
