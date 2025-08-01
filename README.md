# mana

A TUI LLM client for your terminal workflow.

## Why

Using existing LLM clients is a pain. I find myself using `| pbcopy` and pasting into OpenRouter to ask higher level questions that agents can't match. I switch between models for different strengths and test LLM responses without setting up scripts or web UIs. Sometimes I don't want to use an agent, I just want to ask a question and control the context.

Existing command line tools are too CLI-focused or too agentic.

Mana fits into my workflow.

## Features

- [x] Interactive TUI
- [ ] Conversation history
- [ ] Pipe content (`git diff | mana`)
- [ ] Multi-provider support (OpenRouter, Groq, Anthropic, etc.)
- [ ] Model switching
- [ ] Docker image

## Install

```bash
git clone https://github.com/darling/mana.git
cd mana
make build
sudo mv ./bin/mana /usr/local/bin/mana
```

## Setup

Get an API key from [OpenRouter.ai](https://openrouter.ai/) and set:

```bash
export OPENROUTER_API_KEY="your-key-here"
```

More providers coming soon.

## Usage

```bash
mana
```

## Contributing

Fork, branch, commit, PR. Open an issue first for major changes.
