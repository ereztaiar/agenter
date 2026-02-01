package agent

import (
	"context"
	"log"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

type AgentConfig struct {
	Type        AgentType   `yaml:"type"`
	DependsOn   DependsOn   `yaml:"depends_on,omitempty"`
	Name        Name        `yaml:"name"`
	Model       Model       `yaml:"model"`
	Workflow    Workflow    `yaml:"workflow,omitempty"`
	Description Description `yaml:"description"`
	Instruction Instruction `yaml:"instruction"`
	ApiKey      ApiKey      `yaml:"api-key"`
	Arguments   Arguments   `yaml:"arguments,omitempty"`
	Tools       Tools       `yaml:"tools,omitempty"`
	OutputKey   OutputKey   `yaml:"output_key,omitempty"`
	SubAgents   SubAgents   `yaml:"sub_agents,omitempty"`
	Memory      []string    `yaml:"memory,omitempty"`
	agent       *agent.Agent
}

func (ac *AgentConfig) GenerateAgent(agentsConfig map[string]AgentConfig, toolAgents []AgentConfig) {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, string(ac.Model), &genai.ClientConfig{
		APIKey: string(ac.ApiKey),
	})

	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	outputKey := ""
	if ac.OutputKey != nil {
		outputKey = *ac.OutputKey
	}

	var currentAgent agent.Agent

	switch ac.Workflow {
	case "sequential":
		log.Println("building sequential agent")
		agents := ac.SubAgents.BuildAgents(agentsConfig)

		currentAgent, err = sequentialagent.New(sequentialagent.Config{
			AgentConfig: agent.Config{
				Name:        string(ac.Name),
				Description: string(ac.Description),
				SubAgents:   agents,
			},
		})
	case "parallel":
		log.Println("building parallel agent")
		agents := ac.SubAgents.BuildAgents(agentsConfig)

		currentAgent, err = parallelagent.New(parallelagent.Config{
			AgentConfig: agent.Config{
				Name:        string(ac.Name),
				Description: string(ac.Description),
				SubAgents:   agents,
			},
		})
	default:
		var tools []tool.Tool
		if len(toolAgents) > 0 {
			tools = ac.Tools.GenerateAgentTools(toolAgents)
		} else {
			tools = ac.Tools.GenerateFunctionTools()
		}

		currentAgent, err = llmagent.New(llmagent.Config{
			Name:        string(ac.Name),
			Model:       model,
			Description: string(ac.Description),
			Instruction: string(ac.Instruction),
			OutputKey:   outputKey,
			Tools:       tools,
		})
	}

	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	ac.agent = &currentAgent
	log.Println("New", ac.Name, "agent was created")

}

func (ac *AgentConfig) webLauncher(ctx context.Context) {

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(*ac.agent),
	}

	args := strings.Split(string(ac.Arguments), " ")

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, args); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

func (ac *AgentConfig) inlineLauncher(ctx context.Context) {
	log.Println("Starting in-memory command-line runner for agent:", ac.Name)
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(*ac.agent),
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, make([]string, 0)); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}

}
