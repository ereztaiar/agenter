package agent_test

import (
	"testing"

	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRootConfig_Filter(t *testing.T) {
	tests := []struct {
		name       string
		config     agent.RootConfig
		filterType string
		wantCount  int
		wantNames  []string
	}{
		{
			name: "filter root-agent type",
			config: agent.RootConfig{
				AgentsConfig: map[string]agent.AgentConfig{
					"root1": {
						Type: "root-agent",
						Name: "RootAgent1",
					},
					"root2": {
						Type: "root-agent",
						Name: "RootAgent2",
					},
					"agent1": {
						Type: "agent",
						Name: "Agent1",
					},
				},
			},
			filterType: "root-agent",
			wantCount:  2,
			wantNames:  []string{"RootAgent1", "RootAgent2"},
		},
		{
			name: "filter agent type",
			config: agent.RootConfig{
				AgentsConfig: map[string]agent.AgentConfig{
					"root1": {
						Type: "root-agent",
						Name: "RootAgent1",
					},
					"agent1": {
						Type: "agent",
						Name: "Agent1",
					},
					"agent2": {
						Type: "agent",
						Name: "Agent2",
					},
				},
			},
			filterType: "agent",
			wantCount:  2,
			wantNames:  []string{"Agent1", "Agent2"},
		},
		{
			name: "filter non-existent type",
			config: agent.RootConfig{
				AgentsConfig: map[string]agent.AgentConfig{
					"agent1": {
						Type: "agent",
						Name: "Agent1",
					},
				},
			},
			filterType: "tool-agent",
			wantCount:  0,
			wantNames:  []string{},
		},
		{
			name: "filter from empty config",
			config: agent.RootConfig{
				AgentsConfig: map[string]agent.AgentConfig{},
			},
			filterType: "agent",
			wantCount:  0,
			wantNames:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.Filter(tt.filterType)

			assert.Equal(t, tt.wantCount, len(result))

			if tt.wantCount > 0 {
				resultNames := make([]string, len(result))
				for i, ac := range result {
					resultNames[i] = string(ac.Name)
				}

				for _, wantName := range tt.wantNames {
					assert.Contains(t, resultNames, wantName)
				}
			}
		})
	}
}

func TestRootConfig_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		check   func(t *testing.T, config agent.RootConfig)
	}{
		{
			name: "unmarshal simple config with agents",
			yaml: `
agents:
  agent1:
    type: "agent"
    name: "Agent1"
    model: "gemini-2.5-flash-lite"
    description: "Test agent"
    instruction: "Do something"
    api-key: "test-key"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.RootConfig) {
				assert.Len(t, config.AgentsConfig, 1)
				assert.Contains(t, config.AgentsConfig, "agent1")
				assert.Equal(t, agent.AgentType("agent"), config.AgentsConfig["agent1"].Type)
				assert.Equal(t, agent.Name("Agent1"), config.AgentsConfig["agent1"].Name)
				assert.Equal(t, agent.Model("gemini-2.5-flash-lite"), config.AgentsConfig["agent1"].Model)
			},
		},
		{
			name: "unmarshal config with multiple agents",
			yaml: `
agents:
  agent1:
    type: "agent"
    name: "Agent1"
    model: "gemini-2.5-flash-lite"
    description: "Test agent 1"
    instruction: "Do something"
    api-key: "test-key"
  root1:
    type: "root-agent"
    name: "RootAgent"
    model: "gemini-2.5-flash-lite"
    description: "Root agent"
    instruction: "Coordinate"
    api-key: "test-key"
    workflow: "sequential"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.RootConfig) {
				assert.Len(t, config.AgentsConfig, 2)
				assert.Contains(t, config.AgentsConfig, "agent1")
				assert.Contains(t, config.AgentsConfig, "root1")
				assert.Equal(t, agent.AgentType("agent"), config.AgentsConfig["agent1"].Type)
				assert.Equal(t, agent.AgentType("root-agent"), config.AgentsConfig["root1"].Type)
				assert.Equal(t, agent.Workflow("sequential"), config.AgentsConfig["root1"].Workflow)
			},
		},
		{
			name: "unmarshal config with tools",
			yaml: `
agents:
  root1:
    type: "root-agent"
    name: "RootAgent"
    model: "gemini-2.5-flash-lite"
    description: "Root agent"
    instruction: "Coordinate"
    api-key: "test-key"
    tools:
      agents:
        - agent1
        - agent2
      predefined_functions:
        - GoogleSearch
`,
			wantErr: false,
			check: func(t *testing.T, config agent.RootConfig) {
				assert.Len(t, config.AgentsConfig, 1)
				agentConfig := config.AgentsConfig["root1"]
				assert.Equal(t, agent.Agents{"agent1", "agent2"}, agentConfig.Tools.Agents)
				assert.Equal(t, agent.PredefinedFunctions{"GoogleSearch"}, agentConfig.Tools.PredefinedFunctions)
			},
		},
		{
			name: "unmarshal config with memory",
			yaml: `
agents:
  agent1:
    type: "agent"
    name: "Agent1"
    model: "gemini-2.5-flash-lite"
    description: "Test agent"
    instruction: "Do something"
    api-key: "test-key"
memory:
  root-agent-memory: InMemoryRunner
`,
			wantErr: false,
			check: func(t *testing.T, config agent.RootConfig) {
				assert.Len(t, config.MemoryConfig, 1)
				assert.Equal(t, "InMemoryRunner", config.MemoryConfig["root-agent-memory"])
			},
		},
		{
			name: "unmarshal config with output_key",
			yaml: `
agents:
  agent1:
    type: "agent"
    name: "Agent1"
    model: "gemini-2.5-flash-lite"
    description: "Test agent"
    instruction: "Do something"
    api-key: "test-key"
    output_key: "result"
`,
			wantErr: false,
			check: func(t *testing.T, config agent.RootConfig) {
				agent := config.AgentsConfig["agent1"]
				require.NotNil(t, agent.OutputKey)
				assert.Equal(t, "result", *agent.OutputKey)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config agent.RootConfig
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

func TestRootConfig_MarshalYAML(t *testing.T) {
	outputKey := "test_output"
	config := agent.RootConfig{
		AgentsConfig: map[string]agent.AgentConfig{
			"agent1": {
				Type:        "agent",
				Name:        "Agent1",
				Model:       "gemini-2.5-flash-lite",
				Description: "Test agent",
				Instruction: "Do something",
				ApiKey:      "test-key",
				OutputKey:   &outputKey,
			},
			"root1": {
				Type:        "root-agent",
				Name:        "RootAgent",
				Model:       "gemini-2.5-flash-lite",
				Description: "Root agent",
				Instruction: "Coordinate",
				ApiKey:      "test-key",
				Workflow:    "sequential",
				Tools: agent.Tools{
					Agents: agent.Agents{"agent1"},
				},
			},
		},
		MemoryConfig: map[string]string{
			"root-agent-memory": "InMemoryRunner",
		},
	}

	data, err := yaml.Marshal(&config)
	require.NoError(t, err)

	var result agent.RootConfig
	err = yaml.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Len(t, result.AgentsConfig, 2)
	assert.Contains(t, result.AgentsConfig, "agent1")
	assert.Contains(t, result.AgentsConfig, "root1")
	assert.Len(t, result.MemoryConfig, 1)
	assert.Equal(t, "InMemoryRunner", result.MemoryConfig["root-agent-memory"])
}

func TestRootConfig_RoundTrip(t *testing.T) {
	outputKey := "test_output"
	original := agent.RootConfig{
		AgentsConfig: map[string]agent.AgentConfig{
			"agent1": {
				Type:        "agent",
				Name:        "Agent1",
				Model:       "gemini-2.5-flash-lite",
				Description: "Test agent",
				Instruction: "Do something",
				ApiKey:      "test-key",
				OutputKey:   &outputKey,
				Tools: agent.Tools{
					PredefinedFunctions: agent.PredefinedFunctions{"GoogleSearch"},
				},
			},
		},
		MemoryConfig: map[string]string{
			"memory1": "InMemoryRunner",
		},
	}

	data, err := yaml.Marshal(&original)
	require.NoError(t, err)

	var result agent.RootConfig
	err = yaml.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Len(t, result.AgentsConfig, 1)
	assert.Contains(t, result.AgentsConfig, "agent1")

	agent1 := result.AgentsConfig["agent1"]
	assert.Equal(t, agent.AgentType("agent"), agent1.Type)
	assert.Equal(t, agent.Name("Agent1"), agent1.Name)
	require.NotNil(t, agent1.OutputKey)
	assert.Equal(t, "test_output", *agent1.OutputKey)
	assert.Equal(t, agent.PredefinedFunctions{"GoogleSearch"}, agent1.Tools.PredefinedFunctions)

	assert.Len(t, result.MemoryConfig, 1)
	assert.Equal(t, "InMemoryRunner", result.MemoryConfig["memory1"])
}
