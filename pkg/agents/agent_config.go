package agent

import (
	"context"
	"google.golang.org/adk/agent"
	"strings"

	"log"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"

	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
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

func (ac *AgentConfig) GenerateRootAgent(ad []AgentConfig) {
	if ac.Type != "root-agent" {
		log.Fatalln("Wrong agent type")

	}
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, string(ac.Model), &genai.ClientConfig{
		APIKey: string(ac.ApiKey),
	})

	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	var currentAgent agent.Agent

	switch ac.Workflow {
	case "sequential":
		log.Println("building sequential agent")
		agents := []agent.Agent{}

		for _, t := range ad {

			agents = append(agents, *t.agent)
		}

		currentAgent, err = sequentialagent.New(sequentialagent.Config{
			AgentConfig: agent.Config{
				Name:        string(ac.Name),
				Description: string(ac.Description),
				SubAgents:   agents,
			},
		})
	case "parallel":
	default:
		tools := []tool.Tool{}

		for _, t := range ad {
			tool := agenttool.New(*t.agent, nil)
			tools = append(tools, tool)
		}
		currentAgent, err = llmagent.New(llmagent.Config{
			Name:        string(ac.Name),
			Model:       model,
			Description: string(ac.Description),
			Instruction: string(ac.Instruction),
			Tools:       tools,
		})
	}

	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	ac.agent = &currentAgent

	ac.launcher(ctx)

}

func (ac *AgentConfig) launcher(ctx context.Context) {

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(*ac.agent),
	}

	args := strings.Split(string(ac.Arguments), " ")

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, args); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

func (ac *AgentConfig) GenerateAgent() {
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

	tools := []tool.Tool{}

	if ac.Tools.Functions != nil {
		for _, functionTool := range ac.Tools.Functions {
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
		Name:        string(ac.Name),
		Model:       model,
		Description: string(ac.Description),
		Instruction: string(ac.Instruction),
		OutputKey:   outputKey,
		Tools:       tools,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	ac.agent = &currentAgent
	log.Println(`New ` + ac.Name + ` agent was created`)

}
