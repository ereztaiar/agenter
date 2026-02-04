package agent

import (
	"fmt"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/adk/tool/geminitool"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
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

			type tmpArgs struct{}

			type tmpResults struct {
				Symbol string `json:"symbol"`
			}

			executeFnAsync := func(ctx tool.Context, input tmpArgs) (tmpResults, error) {
				log.Println("executing python code")
				cmd := exec.Command("python3", append([]string{functionName}, "")...) //tmpArgs...)...)
				output, err := cmd.Output()
				if err != nil {
					return tmpResults{string("")}, err
				}
				return tmpResults{string(output)}, nil
			}
			var err error
			toolName := strings.TrimSuffix(filepath.Base(functionName), filepath.Ext(functionName))
			functionTool, err = functiontool.New(
				functiontool.Config{
					Name:        toolName,
					Description: fmt.Sprintf("runs %s python method.", toolName),
				},
				executeFnAsync,
			)
			if err != nil {
				log.Fatalln("failed to create long running tool: %w", err)
			}
		}

		tools = append(tools, functionTool)
	}

	return tools
}
