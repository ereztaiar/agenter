package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGetAgentsFile(t *testing.T) {
	tests := []struct {
		name           string
		setupAgentFile func()
		setupViper     func()
		expected       string
	}{
		{
			name: "returns agentFile when set",
			setupAgentFile: func() {
				agentFile = "/custom/path/agent.yaml"
			},
			setupViper: func() {
				viper.Set("file", "/default/path/agent.yaml")
			},
			expected: "/custom/path/agent.yaml",
		},
		{
			name: "returns viper value when agentFile is empty",
			setupAgentFile: func() {
				agentFile = ""
			},
			setupViper: func() {
				viper.Set("file", "/viper/path/agent.yaml")
			},
			expected: "/viper/path/agent.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.setupAgentFile()
			tt.setupViper()

			result := getAgentsFile()
			assert.Equal(t, tt.expected, result)

			agentFile = ""
		})
	}
}

func TestYAMLLoadingWithEnvExpansion(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		envVars map[string]string
		check   func(t *testing.T, config agent.RootConfig)
	}{
		{
			name: "expand environment variables in api-key",
			yaml: `
agents:
  test_agent:
    type: "agent"
    name: "TestAgent"
    model: "gemini-2.5-flash-lite"
    description: "Test"
    instruction: "Test"
    api-key: ${TEST_API_KEY}
`,
			envVars: map[string]string{
				"TEST_API_KEY": "secret-key-123",
			},
			check: func(t *testing.T, config agent.RootConfig) {
				agentConfig := config.AgentsConfig["test_agent"]
				assert.Equal(t, agent.ApiKey("secret-key-123"), agentConfig.ApiKey)
			},
		},
		{
			name: "expand multiple environment variables",
			yaml: `
agents:
  test_agent:
    type: "agent"
    name: "TestAgent"
    model: ${TEST_MODEL}
    description: "Test"
    instruction: "Test"
    api-key: ${TEST_API_KEY}
    arguments: ${TEST_ARGS}
`,
			envVars: map[string]string{
				"TEST_API_KEY": "secret-key-123",
				"TEST_MODEL":   "gemini-2.5-flash-lite",
				"TEST_ARGS":    "web -port 8080",
			},
			check: func(t *testing.T, config agent.RootConfig) {
				agentConfig := config.AgentsConfig["test_agent"]
				assert.Equal(t, agent.ApiKey("secret-key-123"), agentConfig.ApiKey)
				assert.Equal(t, agent.Model("gemini-2.5-flash-lite"), agentConfig.Model)
				assert.Equal(t, agent.Arguments("web -port 8080"), agentConfig.Arguments)
			},
		},
		{
			name: "handle missing environment variables",
			yaml: `
agents:
  test_agent:
    type: "agent"
    name: "TestAgent"
    model: "gemini-2.5-flash-lite"
    description: "Test"
    instruction: "Test"
    api-key: ${MISSING_KEY}
`,
			envVars: map[string]string{},
			check: func(t *testing.T, config agent.RootConfig) {
				agentConfig := config.AgentsConfig["test_agent"]
				assert.Equal(t, agent.ApiKey(""), agentConfig.ApiKey)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			expanded := os.ExpandEnv(tt.yaml)

			var config agent.RootConfig
			err := yaml.Unmarshal([]byte(expanded), &config)
			require.NoError(t, err)

			if tt.check != nil {
				tt.check(t, config)
			}
		})
	}
}

func TestRootCmd_Structure(t *testing.T) {
	assert.NotNil(t, RootCmd)
	assert.Equal(t, "agenter", RootCmd.Use)
	assert.Equal(t, "An AI agent orchestrator", RootCmd.Short)
	assert.NotEmpty(t, RootCmd.Long)
}

func TestListCmd_Structure(t *testing.T) {
	assert.NotNil(t, listCmd)
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List all running agents", listCmd.Short)
	assert.NotEmpty(t, listCmd.Long)
}

func TestRootCmd_HasListSubcommand(t *testing.T) {
	commands := RootCmd.Commands()

	var foundList bool
	for _, cmd := range commands {
		if cmd.Use == "list" {
			foundList = true
			break
		}
	}

	assert.True(t, foundList, "RootCmd should have 'list' subcommand")
}

func TestRootCmd_HasFileFlag(t *testing.T) {
	flag := RootCmd.PersistentFlags().Lookup("file")
	require.NotNil(t, flag, "RootCmd should have --file flag")
	assert.Equal(t, "string", flag.Value.Type())
	assert.Contains(t, flag.Usage, "agent file")
}

func TestViperDefaultFile(t *testing.T) {
	viper.Reset()

	viper.SetDefault("file", filepath.Join(os.Getenv("PWD"), "agent.yaml"))

	defaultFile := viper.GetString("file")
	assert.Contains(t, defaultFile, "agent.yaml")
}

func TestCompleteYAMLLoadingWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	yamlContent := `agents:
  test_agent:
    type: "agent"
    name: "TestAgent"
    model: "gemini-2.5-flash-lite"
    description: "Test agent"
    instruction: "Do something"
    api-key: ${TEST_KEY}
memory:
  test_memory: InMemoryRunner
`

	tmpFile := filepath.Join(tmpDir, "test-agent.yaml")
	err := os.WriteFile(tmpFile, []byte(yamlContent), 0644)
	require.NoError(t, err)

	os.Setenv("TEST_KEY", "my-test-key")
	defer os.Unsetenv("TEST_KEY")

	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	expanded := os.ExpandEnv(string(data))

	var config agent.RootConfig
	err = yaml.Unmarshal([]byte(expanded), &config)
	require.NoError(t, err)

	assert.Len(t, config.AgentsConfig, 1)
	assert.Contains(t, config.AgentsConfig, "test_agent")

	testAgent := config.AgentsConfig["test_agent"]
	assert.Equal(t, agent.AgentType("agent"), testAgent.Type)
	assert.Equal(t, agent.Name("TestAgent"), testAgent.Name)
	assert.Equal(t, agent.ApiKey("my-test-key"), testAgent.ApiKey)

	assert.Len(t, config.MemoryConfig, 1)
	assert.Equal(t, "InMemoryRunner", config.MemoryConfig["test_memory"])
}

func TestYAMLWithComplexTools(t *testing.T) {
	yamlContent := `agents:
  root_agent:
    type: "root-agent"
    name: "RootAgent"
    model: "gemini-2.5-flash-lite"
    description: "Root"
    instruction: "Coordinate"
    api-key: test-key
    workflow: sequential
    tools:
      agents:
        - agent1
        - agent2
      functions:
        - GoogleSearch
  agent1:
    type: "agent"
    name: "Agent1"
    model: "gemini-2.5-flash-lite"
    description: "First agent"
    instruction: "Work"
    api-key: test-key
    output_key: "result1"
  agent2:
    type: "agent"
    name: "Agent2"
    model: "gemini-2.5-flash-lite"
    description: "Second agent uses {result1}"
    instruction: "Process {result1}"
    api-key: test-key
    output_key: "result2"
`

	var config agent.RootConfig
	err := yaml.Unmarshal([]byte(yamlContent), &config)
	require.NoError(t, err)

	assert.Len(t, config.AgentsConfig, 3)

	rootAgent := config.AgentsConfig["root_agent"]
	assert.Equal(t, agent.AgentType("root-agent"), rootAgent.Type)
	assert.Equal(t, agent.Workflow("sequential"), rootAgent.Workflow)
	assert.Len(t, rootAgent.Tools.Agents, 2)
	assert.Contains(t, rootAgent.Tools.Agents, "agent1")
	assert.Contains(t, rootAgent.Tools.Agents, "agent2")
	assert.Len(t, rootAgent.Tools.Functions, 1)
	assert.Contains(t, rootAgent.Tools.Functions, "GoogleSearch")

	agent1 := config.AgentsConfig["agent1"]
	require.NotNil(t, agent1.OutputKey)
	assert.Equal(t, "result1", *agent1.OutputKey)

	agent2 := config.AgentsConfig["agent2"]
	assert.Contains(t, string(agent2.Description), "{result1}")
	assert.Contains(t, string(agent2.Instruction), "{result1}")
	require.NotNil(t, agent2.OutputKey)
	assert.Equal(t, "result2", *agent2.OutputKey)
}
