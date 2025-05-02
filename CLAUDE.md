# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Test Commands
- Build: `go build`
- Run: `go run .` 
- Test: `go test ./...`
- Test specific file: `go test ./path/to/package`
- Test specific function: `go test -run TestFunctionName ./path/to/package`

## Code Style Guidelines
- **Formatting**: Use `gofmt` or `go fmt ./...` to format code
- **Imports**: Group stdlib first, then third-party packages with a blank line between
- **Error handling**: Use direct `if err != nil` checks, return errors with context
- **Naming**: Use camelCase for variables, PascalCase for exported functions/types
- **Comments**: Document exported functions with complete sentences
- **PR conventions**: PR titles should be emojified and descriptive

## Project Structure
- `config/`: Configuration loading and defaults
- `git/`: Git operations and command interfaces
- `llm/`: LLM provider implementations (OpenAI, Anthropic, Gemini)
- `utils/`: Utility functions and helpers