package agent_test

import (
	"testing"
	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)


func TestSubAgents_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    []string
		wantErr bool
	}{
		{
			name: "unmarshal valid agent list",
			yaml: `sub_agents:
  - tech_researcher
  - health_researcher
  - finance_researcher`,
			want:    []string{"tech_researcher", "health_researcher", "finance_researcher"},
			wantErr: false,
		},
		{
			name:    "unmarshal empty list",
			yaml:    `sub_agents: []`,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "unmarshal single agent",
			yaml: `sub_agents:
  - tech_researcher`,
			want:    []string{"tech_researcher"},
			wantErr: false,
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
			assert.Equal(t, len(tt.want), len(config.SubAgents))

			for _, name := range tt.want {
				assert.Contains(t, config.SubAgents, name)
				// Agents are nil until BuildAgents() is called (see sub_agents.go:18)
				assert.Nil(t, config.SubAgents[name])
			}
		})
	}
}

func TestSubAgents_MarshalYAML(t *testing.T) {
	tests := []struct {
		name      string
		subAgents agent.SubAgents
		wantNames []string
		wantNil   bool
	}{
		{
			name: "marshal valid agents",
			subAgents: agent.SubAgents{
				"tech_researcher":    nil,
				"health_researcher":  nil,
				"finance_researcher": nil,
			},
			wantNames: []string{"tech_researcher", "health_researcher", "finance_researcher"},
		},
		{
			name:      "marshal empty map",
			subAgents: agent.SubAgents{},
			wantNil:   true,
		},
		{
			name:      "marshal nil map",
			subAgents: nil,
			wantNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := agent.AgentConfig{
				SubAgents: tt.subAgents,
			}

			data, err := yaml.Marshal(&config)
			require.NoError(t, err)

			if tt.wantNil {
				// Should not contain sub_agents key or it should be empty/null
				assert.NotContains(t, string(data), "sub_agents:")
				return
			}

			// Unmarshal back to verify
			var result agent.AgentConfig
			err = yaml.Unmarshal(data, &result)
			require.NoError(t, err)

			assert.Equal(t, len(tt.wantNames), len(result.SubAgents))
			for _, name := range tt.wantNames {
				assert.Contains(t, result.SubAgents, name)
			}
		})
	}
}

func TestSubAgents_RoundTrip(t *testing.T) {
	original := agent.AgentConfig{
		SubAgents: agent.SubAgents{
			"tech_researcher":   nil,
			"health_researcher": nil,
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&original)
	require.NoError(t, err)

	// Unmarshal back
	var result agent.AgentConfig
	err = yaml.Unmarshal(data, &result)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, len(original.SubAgents), len(result.SubAgents))
	for name := range original.SubAgents {
		assert.Contains(t, result.SubAgents, name)
		// Agents are nil until BuildAgents() is called (see sub_agents.go:18)
		assert.Nil(t, result.SubAgents[name])
	}
}
