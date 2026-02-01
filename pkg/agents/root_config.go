package agent

import "context"

type RootConfig struct {
	AgentsConfig map[string]AgentConfig `yaml:"agents"`
	MemoryConfig map[string]string      `yaml:"memory"`
}

func (rc *RootConfig) BuildAgents() {
	for _, ac := range rc.Filter("root-agent") {

		toolAgentsForRoot := []AgentConfig{}

		for _, agentTool := range ac.Tools.Agents {

			agentConfig := rc.AgentsConfig[agentTool]
			if agentConfig.agent == nil {
				agentConfig.GenerateAgent(rc.AgentsConfig, nil)
			}
			toolAgentsForRoot = append(toolAgentsForRoot, agentConfig)

		}

		ac.GenerateAgent(rc.AgentsConfig, toolAgentsForRoot)

		// Launch the root agent
		ctx := context.Background()
		if ac.Arguments != "" {
			ac.webLauncher(ctx)
		} else {
			ac.inlineLauncher(ctx)
		}
	}
}

// Filter returns all filter agents by string type
func (rc *RootConfig) Filter(t string) []AgentConfig {
	agents := []AgentConfig{}
	for _, agentConfig := range rc.AgentsConfig {
		if agentConfig.Type == AgentType(t) {
			agents = append(agents, agentConfig)
		}
	}
	return agents
}
