# Configuration

There're some sensible defaults configured for the tool.

Check [defaults.go](../config/defaults.go) for mode details.

The config file is stored in `$HOME/.config/gh-aipr/config.yaml`

```yaml
# Force use of a specific provider (optional)
# Valid values: "openai", "anthropic", "gemini"
# If not set, the tool will auto-detect based on available API keys
provider: "anthropic"

# if you choose to override the prompt make sure to check the current one
# in llm/prompt.go as yours going to replace the default one
system_prompt_override: "You are a helpful assistant."

# set the models to ones you want to use
# you don't have to set all of them too
openai:
    model_name: "gpt-4o-mini"

anthropic:
    model_name: "claude-3-5-haiku-latest"

gemini:
    model_name: "gemini-2.5-flash-preview-04-17"
```

## Provider Selection

The tool supports three LLM providers: OpenAI, Anthropic, and Gemini. By default, it will automatically select a provider based on which API keys are available in your environment variables.

If you want to force the use of a specific provider, you can set the `provider` configuration option to one of: `"openai"`, `"anthropic"`, or `"gemini"`. When a specific provider is configured, the tool will:

1. Only attempt to use that provider
2. Show an error if the required API key for that provider is not set
3. Not fall back to other providers

For example, if you set `provider: "anthropic"` but don't have the `ANTHROPIC_API_KEY` environment variable set, the tool will display an error message and exit.
