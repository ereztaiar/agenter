package agent_test

import (
	"testing"

	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestTools_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    agent.Tools
		wantErr bool
	}{
		{
			name: "unmarshal tools with agents only",
			yaml: `tools:
  agents:
    - agent1
    - agent2`,
			want: agent.Tools{
				Agents:    agent.Agents{"agent1", "agent2"},
				PredefinedFunctions: nil,
			},
			wantErr: false,
		},
		{
			name: "unmarshal tools with functions only",
			yaml: `tools:
  predefined_functions:
    - GoogleSearch`,
			want: agent.Tools{
				Agents:    nil,
				PredefinedFunctions: agent.PredefinedFunctions{"GoogleSearch"},
			},
			wantErr: false,
		},
		{
			name: "unmarshal tools with both agents and functions",
			yaml: `tools:
  agents:
    - agent1
  predefined_functions:
    - GoogleSearch`,
			want: agent.Tools{
				Agents:    agent.Agents{"agent1"},
				PredefinedFunctions: agent.PredefinedFunctions{"GoogleSearch"},
			},
			wantErr: false,
		},
		{
			name: "unmarshal empty tools",
			yaml: `tools: {}`,
			want: agent.Tools{
				Agents:    nil,
				PredefinedFunctions: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config struct {
				Tools agent.Tools `yaml:"tools"`
			}
			err := yaml.Unmarshal([]byte(tt.yaml), &config)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.Agents, config.Tools.Agents)
			assert.Equal(t, tt.want.PredefinedFunctions, config.Tools.PredefinedFunctions)
		})
	}
}

func TestTools_MarshalYAML(t *testing.T) {
	tests := []struct {
		name  string
		tools agent.Tools
	}{
		{
			name: "marshal tools with agents and functions",
			tools: agent.Tools{
				Agents:    agent.Agents{"agent1", "agent2"},
				PredefinedFunctions: agent.PredefinedFunctions{"GoogleSearch"},
			},
		},
		{
			name: "marshal tools with agents only",
			tools: agent.Tools{
				Agents:    agent.Agents{"agent1"},
				PredefinedFunctions: nil,
			},
		},
		{
			name: "marshal tools with functions only",
			tools: agent.Tools{
				Agents:    nil,
				PredefinedFunctions: agent.PredefinedFunctions{"GoogleSearch"},
			},
		},
		{
			name: "marshal empty tools",
			tools: agent.Tools{
				Agents:    nil,
				PredefinedFunctions: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := struct {
				Tools agent.Tools `yaml:"tools,omitempty"`
			}{
				Tools: tt.tools,
			}

			data, err := yaml.Marshal(&config)
			require.NoError(t, err)

			var result struct {
				Tools agent.Tools `yaml:"tools"`
			}
			err = yaml.Unmarshal(data, &result)
			require.NoError(t, err)

			assert.Equal(t, tt.tools.Agents, result.Tools.Agents)
			assert.Equal(t, tt.tools.PredefinedFunctions, result.Tools.PredefinedFunctions)
		})
	}
}

func TestTools_RoundTrip(t *testing.T) {
	original := agent.Tools{
		Agents:    agent.Agents{"agent1", "agent2", "agent3"},
		PredefinedFunctions: agent.PredefinedFunctions{"GoogleSearch", "CustomFunction"},
	}

	config := struct {
		Tools agent.Tools `yaml:"tools"`
	}{
		Tools: original,
	}

	data, err := yaml.Marshal(&config)
	require.NoError(t, err)

	var result struct {
		Tools agent.Tools `yaml:"tools"`
	}
	err = yaml.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, original.Agents, result.Tools.Agents)
	assert.Equal(t, original.PredefinedFunctions, result.Tools.PredefinedFunctions)
}
