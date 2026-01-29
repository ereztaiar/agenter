package agent

import (
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
)

type Tools struct {
	Agents    Agents    `yaml:"agents,omitempty"`
	Functions Functions `yaml:"functions,omitempty"`
}

// GenerateAgentTools converts a list of AgentConfig into agent tools
func (t *Tools) GenerateAgentTools(agents []AgentConfig) []tool.Tool {
	tools := []tool.Tool{}

	for _, agentConfig := range agents {
		if agentConfig.agent != nil {
			agentTool := agenttool.New(*agentConfig.agent, nil)
			tools = append(tools, agentTool)
		}
	}

	return tools
}

// GenerateFunctionTools converts the Functions list into function tools
func (t *Tools) GenerateFunctionTools() []tool.Tool {
	tools := []tool.Tool{}

	if t.Functions == nil {
		return tools
	}

	for _, functionName := range t.Functions {
		var functionTool tool.Tool

		switch functionName {
		case "GoogleSearch":
			functionTool = geminitool.GoogleSearch{}
		default:
			// Skip unknown function tools
			continue
		}

		tools = append(tools, functionTool)
	}

	return tools
}