# Configuration

There're some sensible defaults configured for the tool.

Check [defaults.go](../config/defaults.go) for mode details.

The config file is stored in `$HOME/.config/gh-aipr/config.yaml`

```yaml
# Completely replaces the default system prompt
# in llm/prompt.go with your custom prompt
system_prompt_override: "You are a helpful assistant."

# Extends the system prompt (or override) with additional instructions
# These will be appended to the prompt with a clear section header
additional_instructions: "Focus on security implications of the changes."

# set the models to ones you want to use
# you don't have to set all of them too
openai:
    model_name: "gpt-4o-mini"

anthropic:
    model_name: "claude-3-5-haiku-latest"

gemini:
    model_name: "gemini-2.5-flash-preview-04-17"
```
