package agent_test

import (
	"testing"

	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestAgentConfig_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		check   func(t *testing.T, config agent.AgentConfig)
	}{
		{
			name: "unmarshal basic agent config",
			yaml: `
type: "agent"
name: "TestAgent"
model: "gemini-2.5-flash-lite"
description: "A test agent"
instruction: "Do something useful"
api-key: "test-api-key"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				assert.Equal(t, agent.AgentType("agent"), config.Type)
				assert.Equal(t, agent.Name("TestAgent"), config.Name)
				assert.Equal(t, agent.Model("gemini-2.5-flash-lite"), config.Model)
				assert.Equal(t, agent.Description("A test agent"), config.Description)
				assert.Equal(t, agent.Instruction("Do something useful"), config.Instruction)
				assert.Equal(t, agent.ApiKey("test-api-key"), config.ApiKey)
			},
		},
		{
			name: "unmarshal root agent with workflow",
			yaml: `
type: "root-agent"
name: "RootAgent"
model: "gemini-2.5-flash-lite"
description: "Root orchestrator"
instruction: "Coordinate agents"
api-key: "test-api-key"
workflow: "sequential"
arguments: "web -port 8080"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				assert.Equal(t, agent.AgentType("root-agent"), config.Type)
				assert.Equal(t, agent.Workflow("sequential"), config.Workflow)
				assert.Equal(t, agent.Arguments("web -port 8080"), config.Arguments)
			},
		},
		{
			name: "unmarshal agent with output_key",
			yaml: `
type: "agent"
name: "ResearchAgent"
model: "gemini-2.5-flash-lite"
description: "Research agent"
instruction: "Research topics"
api-key: "test-api-key"
output_key: "research_results"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				require.NotNil(t, config.OutputKey)
				assert.Equal(t, "research_results", *config.OutputKey)
			},
		},
		{
			name: "unmarshal agent without output_key",
			yaml: `
type: "agent"
name: "SimpleAgent"
model: "gemini-2.5-flash-lite"
description: "Simple agent"
instruction: "Do work"
api-key: "test-api-key"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				assert.Nil(t, config.OutputKey)
			},
		},
		{
			name: "unmarshal agent with tools",
			yaml: `
type: "agent"
name: "ToolAgent"
model: "gemini-2.5-flash-lite"
description: "Agent with tools"
instruction: "Use tools"
api-key: "test-api-key"
tools:
  predefined_functions:
    - GoogleSearch
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				assert.Equal(t, agent.PredefinedFunctions{"GoogleSearch"}, config.Tools.PredefinedFunctions)
			},
		},
		{
			name: "unmarshal agent with sub_agents",
			yaml: `
type: "root-agent"
name: "MultiAgent"
model: "gemini-2.5-flash-lite"
description: "Multi-agent system"
instruction: "Coordinate"
api-key: "test-api-key"
sub_agents:
  - agent1
  - agent2
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				assert.Len(t, config.SubAgents, 2)
				assert.Contains(t, config.SubAgents, "agent1")
				assert.Contains(t, config.SubAgents, "agent2")
			},
		},
		{
			name: "unmarshal agent with depends_on",
			yaml: `
type: "agent"
name: "DependentAgent"
model: "gemini-2.5-flash-lite"
description: "Dependent agent"
instruction: "Work after other"
api-key: "test-api-key"
depends_on: "other_agent"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				assert.Equal(t, agent.DependsOn("other_agent"), config.DependsOn)
			},
		},
		{
			name: "unmarshal agent with memory",
			yaml: `
type: "agent"
name: "MemoryAgent"
model: "gemini-2.5-flash-lite"
description: "Agent with memory"
instruction: "Remember things"
api-key: "test-api-key"
memory:
  - memory1
  - memory2
`,
			wantErr: false,
			check: func(t *testing.T, config agent.AgentConfig) {
				assert.Equal(t, []string{"memory1", "memory2"}, config.Memory)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config agent.AgentConfig
			err := yaml.Unmarshal([]byte(tt.yaml), &config)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, config)
			}
		})
	}
}

func TestAgentConfig_MarshalYAML(t *testing.T) {
	tests := []struct {
		name   string
		config agent.AgentConfig
	}{
		{
			name: "marshal basic agent",
			config: agent.AgentConfig{
				Type:        "agent",
				Name:        "TestAgent",
				Model:       "gemini-2.5-flash-lite",
				Description: "Test agent",
				Instruction: "Do work",
				ApiKey:      "test-key",
			},
		},
		{
			name: "marshal root agent with workflow",
			config: agent.AgentConfig{
				Type:        "root-agent",
				Name:        "RootAgent",
				Model:       "gemini-2.5-flash-lite",
				Description: "Root",
				Instruction: "Coordinate",
				ApiKey:      "test-key",
				Workflow:    "sequential",
				Arguments:   "web -port 8080",
			},
		},
		{
			name: "marshal agent with output_key",
			config: agent.AgentConfig{
				Type:        "agent",
				Name:        "ResearchAgent",
				Model:       "gemini-2.5-flash-lite",
				Description: "Research",
				Instruction: "Research",
				ApiKey:      "test-key",
				OutputKey:   func() *string { s := "results"; return &s }(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(&tt.config)
			require.NoError(t, err)

			var result agent.AgentConfig
			err = yaml.Unmarshal(data, &result)
			require.NoError(t, err)

			assert.Equal(t, tt.config.Type, result.Type)
			assert.Equal(t, tt.config.Name, result.Name)
			assert.Equal(t, tt.config.Model, result.Model)
			assert.Equal(t, tt.config.Description, result.Description)
			assert.Equal(t, tt.config.Instruction, result.Instruction)
		})
	}
}

func TestAgentConfig_RoundTrip(t *testing.T) {
	outputKey := "test_output"
	original := agent.AgentConfig{
		Type:        "agent",
		Name:        "TestAgent",
		Model:       "gemini-2.5-flash-lite",
		Description: "Test agent",
		Instruction: "Do something",
		ApiKey:      "test-api-key",
		OutputKey:   &outputKey,
		Tools: agent.Tools{
			PredefinedFunctions: agent.PredefinedFunctions{"GoogleSearch"},
		},
		Memory: []string{"memory1"},
	}

	data, err := yaml.Marshal(&original)
	require.NoError(t, err)

	var result agent.AgentConfig
	err = yaml.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, original.Type, result.Type)
	assert.Equal(t, original.Name, result.Name)
	assert.Equal(t, original.Model, result.Model)
	assert.Equal(t, original.Description, result.Description)
	assert.Equal(t, original.Instruction, result.Instruction)
	assert.Equal(t, original.ApiKey, result.ApiKey)
	require.NotNil(t, result.OutputKey)
	assert.Equal(t, *original.OutputKey, *result.OutputKey)
	assert.Equal(t, original.Tools.PredefinedFunctions, result.Tools.PredefinedFunctions)
	assert.Equal(t, original.Memory, result.Memory)
}

func TestAgentConfig_ComplexConfiguration(t *testing.T) {
	yamlConfig := `
type: "root-agent"
name: "ComplexAgent"
model: "gemini-2.5-flash-lite"
description: "Complex multi-agent system"
instruction: "Coordinate multiple agents with tools"
api-key: "test-api-key"
workflow: "parallel"
arguments: "web -port 8081 api webui"
tools:
  agents:
    - researcher
    - writer
    - editor
  predefined_functions:
    - GoogleSearch
sub_agents:
  - sub1
  - sub2
memory:
  - shared_memory
output_key: "final_result"
depends_on: "prerequisite_agent"
`

	var config agent.AgentConfig
	err := yaml.Unmarshal([]byte(yamlConfig), &config)
	require.NoError(t, err)

	assert.Equal(t, agent.AgentType("root-agent"), config.Type)
	assert.Equal(t, agent.Name("ComplexAgent"), config.Name)
	assert.Equal(t, agent.Model("gemini-2.5-flash-lite"), config.Model)
	assert.Equal(t, agent.Workflow("parallel"), config.Workflow)
	assert.Equal(t, agent.Arguments("web -port 8081 api webui"), config.Arguments)
	assert.Equal(t, agent.DependsOn("prerequisite_agent"), config.DependsOn)

	assert.Len(t, config.Tools.Agents, 3)
	assert.Contains(t, config.Tools.Agents, "researcher")
	assert.Contains(t, config.Tools.Agents, "writer")
	assert.Contains(t, config.Tools.Agents, "editor")

	assert.Len(t, config.Tools.PredefinedFunctions, 1)
	assert.Contains(t, config.Tools.PredefinedFunctions, "GoogleSearch")

	assert.Len(t, config.SubAgents, 2)
	assert.Contains(t, config.SubAgents, "sub1")
	assert.Contains(t, config.SubAgents, "sub2")

	assert.Equal(t, []string{"shared_memory"}, config.Memory)

	require.NotNil(t, config.OutputKey)
	assert.Equal(t, "final_result", *config.OutputKey)
}

func TestAgentType_Values(t *testing.T) {
	tests := []struct {
		name      string
		agentType agent.AgentType
		expected  string
	}{
		{"agent type", "agent", "agent"},
		{"root-agent type", "root-agent", "root-agent"},
		{"tool-agent type", "tool-agent", "tool-agent"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.agentType))
		})
	}
}

func TestWorkflow_Values(t *testing.T) {
	tests := []struct {
		name     string
		workflow agent.Workflow
		expected string
	}{
		{"sequential workflow", "sequential", "sequential"},
		{"parallel workflow", "parallel", "parallel"},
		{"empty workflow", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.workflow))
		})
	}
}
