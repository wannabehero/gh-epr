# Configuration

There're some sensible defaults configured for the tool.

Check [defaults.go](../config/defaults.go) for mode details.

## Global Configuration

The global config file is stored in `$HOME/.config/gh-aipr/config.yaml`

```yaml
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

## Local Repository Configuration

You can also create a local configuration file within a git repository that will be merged with the global configuration. 

Local configuration should be placed at `.github/gh-aipr/config.yaml` in your repository.

This might be useful if you want to specify model or instructions on a per-repository basis.

The local configuration will override the global configuration for identical keys, and both configurations will be deep merged.
