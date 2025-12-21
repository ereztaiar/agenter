package agent

import (
	"context"
	"strings"

	"log"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"

	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
	// "google.golang.org/adk/memory"
)

func (rc *RootConfig) BuildAgents() {
	rc.buildToolAgents()
	rc.buildRootAgents()
}

func (rc *RootConfig) buildToolAgents() {
	for _, ac := range rc.Filter("tool-agent") {
		ac.GenerateAgent()
	}
}

func (rc *RootConfig) buildRootAgents() {
	for _, ac := range rc.Filter("root-agent") {

		toolAgentsForRoot := []AgentConfig{}
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
		if agentConfig.Type == t {
			agents = append(agents, agentConfig)
		}
	}
	return agents
}

func (ac *AgentConfig) GenerateRootAgent(ad []AgentConfig) {
	if ac.Type != "root-agent" {
		log.Fatalln("Wrong agent type")

	}
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, ac.Model, &genai.ClientConfig{
		APIKey: ac.ApiKey,
	})

	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	tools := []tool.Tool{}

	for _, t := range ad {
		tool := agenttool.New(*t.agent, nil)
		tools = append(tools, tool)
	}

	currentAgent, err := llmagent.New(llmagent.Config{
		Name:        ac.Name,
		Model:       model,
		Description: ac.Description,
		Instruction: ac.Instruction,
		Tools:       tools,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	ac.agent = &currentAgent

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(currentAgent),
	}

	args := strings.Split(ac.Arguments, " ")

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, args); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}

}

func (ac *AgentConfig) GenerateAgent() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, ac.Model, &genai.ClientConfig{
		APIKey: ac.ApiKey,
	})

	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	outputKey := ""
	if ac.OutputKey != nil {
		outputKey = *ac.OutputKey
	}

	tools := []tool.Tool{}


	if ac.Tools.Function != nil {
		for _, functionTool := range ac.Tools.Function {
			var t tool.Tool
			switch functionTool {
			case "":
				t = geminitool.GoogleSearch{}
			default:
				continue
			}

			tools = append(tools, t)

		}
	}

	currentAgent, err := llmagent.New(llmagent.Config{
		Name:        ac.Name,
		Model:       model,
		Description: ac.Description,
		Instruction: ac.Instruction,
		OutputKey:   outputKey,
		Tools:       tools,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	ac.agent = &currentAgent
	log.Println(`New ` + ac.Name + ` agent was created`)

}
