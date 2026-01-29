# Agenter

A Go-based AI agent orchestrator that uses Google's Agent Development Kit (ADK) to build and run multi-agent systems. Agenter enables you to define complex AI workflows using simple YAML configuration files, where agents can be composed into sequential pipelines, parallel executions, or LLM-driven orchestrations.

## Features

- **YAML-Based Configuration**: Define your entire agent workflow in declarative YAML files
- **Multiple Workflow Modes**:
  - **Sequential**: Execute agents in a predefined order
  - **Parallel**: Run multiple agents concurrently
  - **LLM-Driven**: Let the AI decide which agents to call based on the task
- **Agent Composition**: Build complex workflows by composing agents as tools for other agents
- **Output Chaining**: Pass outputs from one agent to another using template variables
- **Function Tools**: Integrate built-in tools like Google Search into your agents
- **Web UI & API**: Built-in web interface and REST API for interacting with agents
- **Environment Variable Support**: Dynamic configuration with environment variable substitution

## Installation

### Prerequisites

- Go 1.24.4 or later
- Google API key for Gemini models

### Build from Source

```bash
# Clone the repository
git clone https://github.com/ereztaiar/agenter.git
cd agenter

# Build the binary
go build -o agenter ./cmd/agent
```

## Quick Start

### 1. Set Up Environment

Create a `.env` file in the root directory:

```env
GOOGLE_API_KEY=your_google_api_key_here
ROOT_AGENT_ARGUMENTS=--web --port 8080
```

### 2. Create Your First Agent Configuration

Create an `agent.yaml` file:

```yaml
agents:
  writer_agent:
    type: "agent"
    name: "WriterAgent"
    model: "gemini-2.5-flash-lite"
    description: "Writes content based on user input."
    instruction: >
      Write a short, engaging paragraph about the given topic.
      Keep it concise and informative.
    api-key: ${GOOGLE_API_KEY}
    output_key: "content"

  root_agent:
    type: "root-agent"
    name: "ContentCreator"
    model: "gemini-2.5-flash-lite"
    arguments: ${ROOT_AGENT_ARGUMENTS}
    tools:
      agents:
        - writer_agent
```

### 3. Run the Agent

```bash
# Run with default agent.yaml file
./agenter

# Run with specific configuration file
./agenter --file agent.yaml
```

The agent will start a web server (if configured) and be ready to accept requests.

## Configuration Format

### Agent Types

#### Leaf Agent (type: "agent")

Basic AI agent that performs a specific task.

```yaml
agent_name:
  type: "agent"
  name: "DisplayName"
  model: "gemini-2.5-flash-lite"
  description: "What this agent does"
  instruction: >
    Detailed instructions for the agent.
    Use {output_key} to reference other agents' outputs.
  api-key: ${GOOGLE_API_KEY}
  output_key: "result_key"  # Makes output available to other agents
  tools:
    functions:
      - GoogleSearch  # Add function tools
```

#### Root Agent (type: "root-agent")

Orchestration agent that coordinates other agents.

```yaml
root_agent:
  type: "root-agent"
  name: "Orchestrator"
  model: "gemini-2.5-flash-lite"
  arguments: ${ROOT_AGENT_ARGUMENTS}  # e.g., "--web --port 8080"
  workflow: "sequential"  # Optional: "sequential", "parallel", or omit for LLM-driven
  tools:
    agents:
      - agent1
      - agent2
```

### Workflow Modes

1. **Sequential Workflow**: Agents execute in the order specified
   ```yaml
   workflow: "sequential"
   ```

2. **Parallel Workflow**: Agents execute concurrently
   ```yaml
   workflow: "parallel"
   ```

3. **LLM-Driven** (default): The LLM decides which agents to call based on instructions

### Output Chaining

Pass data between agents using output keys:

```yaml
agents:
  agent1:
    type: "agent"
    instruction: "Analyze the topic and create an outline."
    output_key: "outline"

  agent2:
    type: "agent"
    instruction: "Write content based on this outline: {outline}"
    # References agent1's output
```

## Examples

### Example 1: Sequential Blog Writing Pipeline

```yaml
agents:
  outline_agent:
    type: "agent"
    name: "OutlineAgent"
    model: "gemini-2.5-flash-lite"
    description: "Creates the initial blog post outline."
    instruction: >
      Create a blog outline with:
      1. A catchy headline
      2. An introduction hook
      3. 3-5 main sections
      4. A concluding thought
    api-key: ${GOOGLE_API_KEY}
    output_key: "blog_outline"

  writer_agent:
    type: "agent"
    name: "WriterAgent"
    model: "gemini-2.5-flash-lite"
    description: "Writes the full blog post."
    instruction: >
      Following this outline: {blog_outline}
      Write a 200-300 word blog post.
    api-key: ${GOOGLE_API_KEY}
    output_key: "blog_draft"

  editor_agent:
    type: "agent"
    name: "EditorAgent"
    model: "gemini-2.5-flash-lite"
    description: "Edits and polishes the draft."
    instruction: >
      Edit this draft: {blog_draft}
      Fix grammatical errors and improve flow.
    api-key: ${GOOGLE_API_KEY}
    output_key: "final_blog"

  root_agent:
    type: "root-agent"
    workflow: "sequential"
    name: "BlogPipeline"
    model: "gemini-2.5-flash-lite"
    arguments: ${ROOT_AGENT_ARGUMENTS}
    tools:
      agents:
        - outline_agent
        - writer_agent
        - editor_agent
```

### Example 2: Research Agent with Google Search

```yaml
agents:
  research_agent:
    type: "agent"
    name: "ResearchAgent"
    model: "gemini-2.5-flash-lite"
    description: "Research agent with web search capability."
    instruction: >
      Use google_search to find 2-3 pieces of relevant information
      on the given topic and present findings with citations.
    api-key: ${GOOGLE_API_KEY}
    tools:
      functions:
        - GoogleSearch
    output_key: "research_findings"

  summarizer_agent:
    type: "agent"
    name: "SummarizerAgent"
    model: "gemini-2.5-flash-lite"
    description: "Summarizes research findings."
    instruction: >
      Read the research findings: {research_findings}
      Create a concise summary as a bulleted list.
    api-key: ${GOOGLE_API_KEY}
    output_key: "final_summary"

  root_agent:
    type: "root-agent"
    name: "ResearchCoordinator"
    model: "gemini-2.5-flash-lite"
    arguments: ${ROOT_AGENT_ARGUMENTS}
    tools:
      agents:
        - research_agent
        - summarizer_agent
```

## Architecture

### Core Components

- **RootConfig**: Top-level configuration that loads YAML files and orchestrates agent building
- **AgentConfig**: Individual agent configuration with type, model, instructions, and tools
- **Tools System**: Manages both agent tools (sub-agents) and function tools (like GoogleSearch)
- **Launcher**: Google ADK's launcher provides web UI, API server, and CLI interfaces

### Configuration Flow

1. Load YAML file and perform environment variable substitution
2. Unmarshal into RootConfig structure
3. Build agent dependency tree
4. Generate root agents with their tool agents
5. Launch with Google ADK's launcher framework

### File Structure

```
agenter/
├── cmd/agent/           # Application entry point
│   └── main.go
├── pkg/
│   ├── agents/          # Agent configuration and building
│   │   ├── root_config.go
│   │   ├── agent_config.go
│   │   ├── agent.go
│   │   ├── tools.go
│   │   └── sub_agents.go
│   └── commands/        # CLI commands
│       └── commands.go
├── agent.yaml           # Example configurations
├── .env                 # Environment variables
└── README.md
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test
go test ./pkg/agents -run TestSubAgents_UnmarshalYAML

# Run tests with verbose output
go test -v ./...
```

### Available Commands

```bash
# Start the agent
./agenter

# Use specific configuration file
./agenter --file agent.second.yaml

# List agents (not yet implemented)
./agenter list
```

## Configuration Reference

### Supported Models

- `gemini-2.5-flash-lite`
- `gemini-2.5-flash`
- `gemini-2.5-pro`
- Other Gemini models supported by Google's Genai SDK

### Supported Function Tools

- `GoogleSearch`: Web search capability

### Launch Arguments

Configure the root agent's runtime behavior with `arguments`:

- `--web`: Enable web UI
- `--port <number>`: Specify port (default: 8080)
- `--cli`: Enable CLI interface
- `--api`: Enable REST API

Example:
```yaml
arguments: "--web --port 8080"
```

## Environment Variables

- `GOOGLE_API_KEY`: Required. Your Google Gemini API key
- `ROOT_AGENT_ARGUMENTS`: Optional. Default launch arguments for root agents

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

By contributing to this project, you agree to the terms outlined in our [Contributor License Agreement (CLA)](CLA.md). Simply submitting a pull request constitutes your acceptance of the CLA terms.

## Support

For issues and questions, please open an issue on GitHub.
