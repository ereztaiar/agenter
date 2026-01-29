package agent

type RootConfig struct {
	AgentsConfig map[string]AgentConfig `yaml:"agents"`
	MemoryConfig map[string]string      `yaml:"memory"`
}

func (rc *RootConfig) BuildAgents() {
	// rc.buildToolAgents()
	rc.buildRootAgents()
}

func (rc *RootConfig) buildRootAgents() {
	for _, ac := range rc.Filter("root-agent") {

		toolAgentsForRoot := []AgentConfig{}

		ac.SubAgents.buildSubagents()

		for _, agentTool := range ac.Tools.Agents {

			agentConfig := rc.AgentsConfig[agentTool]
			if agentConfig.agent == nil {
				agentConfig.GenerateAgent()
			}
			toolAgentsForRoot = append(toolAgentsForRoot, agentConfig)

		}

		ac.GenerateRootAgent(toolAgentsForRoot)
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
