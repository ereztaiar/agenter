package agent

import (
	"google.golang.org/adk/agent"
)

type Tools struct {
	Agents   []string `yaml:"agents,omitempty"`
	Function []string `yaml:"functions,omitempty"`
}

type AgentConfig struct {
	Type        string   `yaml:"type"`
	DependsOn   string   `yaml:"depends_on,omitempty"`
	Name        string   `yaml:"name"`
	Model       string   `yaml:"model"`
	Workflow    string   `yaml:"workflow,omitempty"`
	Description string   `yaml:"description"`
	Instruction string   `yaml:"instruction"`
	ApiKey      string   `yaml:"api-key"`
	Arguments   string   `yaml:"arguments,omitempty"`
	Tools       Tools    `yaml:"tools,omitempty"`
	OutputKey   *string  `yaml:"output_key,omitempty"`
	Memory      []string `yaml:"memory,omitempty"`
	agent       *agent.Agent
}

type RootConfig struct {
	AgentsConfig map[string]AgentConfig `yaml:"agents"`
	MemoryConfig map[string]string      `yaml:"memory"`
}
