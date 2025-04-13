# gh-epr

Github CLI extension to add emojis to pull requests titles.

## Installation

Make sure the CLI is installed: https://cli.github.com/

```
gh extension install wannabehero/gh-epr
```

## Configuration

To make the best use of the tool set one of the following:
- `OPENAI_API_KEY`
- `ANTHROPIC_API_KEY`
- `GEMINI_API_KEY`
in your environment variables so it can use the LLM
to generate relevant title automatically
based on your current branch commits.

## Usage

```
gh epr <other args that gh pr create supports>
```

This will follow you through interactive create process
and your final PR title will be automatically added and emojified.
