# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Agenter is a Go-based AI agent orchestrator that uses Google's ADK (Agent Development Kit) to build and run multi-agent systems. It enables YAML-based configuration of agent workflows, where agents can be composed into sequential pipelines or parallel executions.

## Build & Run Commands

```bash
# Build the binary
go build -o agenter ./cmd/agent

# Run with default agent.yaml file
./agenter

# Run with specific agent file
./agenter --file agent.second.yaml

# List agents (command exists but not yet implemented)
./agenter list

# Run tests
go test ./...

# Run specific test
go test ./pkg/agents -run TestSubAgents_UnmarshalYAML
```

## Environment Setup

The application requires a `.env` file in the root directory with:
- `GOOGLE_API_KEY`: API key for Google Gemini models
- `ROOT_AGENT_ARGUMENTS`: Runtime arguments for root agent (e.g., web server config)

## Architecture

### Core Components

**Agent Configuration System** (`pkg/agents/`)
- `RootConfig`: Top-level config that loads YAML files and orchestrates agent building
- `AgentConfig`: Individual agent configuration with type, model, instructions, tools
- Agent types: `agent` (leaf agents), `root-agent` (orchestrators)
- Workflow modes: `sequential` (ordered execution), `parallel` (concurrent), or default (LLM-driven tool selection)

**Configuration Flow**
1. `commands.loadAgents()` reads YAML file and performs `os.ExpandEnv()` for variable substitution
2. YAML unmarshals into `RootConfig` structure
3. `RootConfig.BuildAgents()` filters for root-agents and builds dependency tree
4. Each root-agent generates its tool agents from the `tools.agents` list
5. Root agent launches with Google ADK's launcher framework

**Agent Types**
- **Leaf Agents**: Basic LLM agents with model, instructions, and optional tools (functions like GoogleSearch)
- **Root Agents**: Orchestration agents that wrap other agents as tools
  - Sequential workflow: Executes agents in order using `sequentialagent.New()`
  - Default workflow: LLM chooses which agent-tools to call using `llmagent.New()`

**Tools System** (`pkg/agents/tools.go`)
- `Tools.Agents`: List of agent names to use as sub-agents (converted to `agenttool.New()`)
- `Tools.Functions`: List of function tool names (e.g., "GoogleSearch" → `geminitool.GoogleSearch{}`)

**Key Mechanisms**
- Agents communicate through `output_key`: one agent's output becomes available to downstream agents via `{output_key}` template syntax in instructions
- The `SubAgents` type uses custom YAML marshal/unmarshal to store agent references as string names in YAML but as `*agent.Agent` in memory
- Launcher uses Google ADK's `full.NewLauncher()` which provides web UI, API server, and CLI interfaces based on `ROOT_AGENT_ARGUMENTS`

### File Structure

```
cmd/agent/main.go          # Entry point, calls commands.Start()
pkg/commands/commands.go   # Cobra CLI setup, YAML loading, agent lifecycle
pkg/agents/
  ├── root_config.go       # Top-level config, agent building orchestration
  ├── agent_config.go      # Agent generation (root & leaf), launcher setup
  ├── agent.go             # Type definitions (Name, Model, Workflow, etc.)
  ├── tools.go             # Tools structure (Agents, Functions)
  └── sub_agents.go        # Custom YAML marshaling for agent references
```

### YAML Configuration Format

```yaml
agents:
  agent_name:
    type: "agent"              # or "root-agent"
    name: "DisplayName"
    model: "gemini-2.5-flash-lite"
    description: "What this agent does"
    instruction: >
      Instructions for the agent.
      Use {output_key} to reference other agents' outputs.
    api-key: ${GOOGLE_API_KEY}
    output_key: "result_key"   # Optional: makes output available to other agents
    tools:
      agents:                  # Sub-agents as tools (root-agent only)
        - sub_agent_1
        - sub_agent_2
      functions:               # Function tools (currently only GoogleSearch)
        - GoogleSearch
    workflow: "sequential"     # Optional: "sequential" or "parallel" (root-agent only)
    arguments: "..."           # Optional: launcher args (root-agent only)
```

## Development Notes

- The codebase uses Google's ADK (`google.golang.org/adk`) and Genai (`google.golang.org/genai`) packages
- Cobra for CLI, Viper for config management, godotenv for env loading
- Test files use testify for assertions
- The `SubAgents.buildSubagents()` method exists but is currently a no-op placeholder
- Memory configuration is defined in YAML but not yet implemented in the agent generation code
