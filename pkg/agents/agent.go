package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	// "google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
	// "google.golang.org/adk/memory"
)

type AgentConfig struct {
	Type        string   `yaml:"type"`
	DependsOn   string   `yaml:"depends_on,omitempty"`
	Name        string   `yaml:"name"`
	Model       string   `yaml:"model"`
	Description string   `yaml:"description"`
	Instruction string   `yaml:"instruction"`
	ApiKey      string   `yaml:"api-key"`
	Arguments   string   `yaml:"arguments,omitempty"`
	Tools       []string `yaml:"tools,omitempty"`
	Memory      []string `yaml:"memory,omitempty"`
	agent       *agent.Agent
}

type RootConfig struct {
	AgentsConfig map[string]AgentConfig `yaml:"agents"`
	MemoryConfig map[string]string      `yaml:"memory"`
}

// GetAgentsDependents returns all dependent agents by string name
func (rc *RootConfig) GetAgentsDependents(n string) []agent.Agent {
	dependents := []agent.Agent{}
	for _, agent := range rc.AgentsConfig {
		if agent.DependsOn == n {
			dependents = append(dependents, *agent.agent)
		}
	}
	return dependents
}

// Filter returns all filter agents by string type
func (rc *RootConfig) Filter(t string) []agent.Agent {
	agents := []agent.Agent{}
	for _, agent := range rc.AgentsConfig {
		if agent.Type == t {
			agents = append(agents, *agent.agent)
		}
	}
	return agents
}

func (ac *AgentConfig) GenerateAgent() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, ac.Model, &genai.ClientConfig{
		APIKey: ac.ApiKey,
	})

	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	switch ac.Type {
	case "standalone-agent":
		currentAgent, err := llmagent.New(llmagent.Config{
			Name:        ac.Name,
			Model:       model,
			Description: ac.Description,
			Instruction: ac.Instruction,
			Tools: []tool.Tool{
				geminitool.GoogleSearch{},
			},
		})
		if err != nil {
			log.Fatalf("Failed to create agent: %v", err)
		}

		config := &launcher.Config{
			AgentLoader: agent.NewSingleLoader(currentAgent),
		}

		l := full.NewLauncher()
		args := strings.Split(ac.Arguments, " ")

		if err = l.Execute(ctx, config, args); err != nil {
			log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
		}

	case "tool-agent":
		currentAgent, err := llmagent.New(llmagent.Config{
			Name:        ac.Name,
			Model:       model,
			Description: ac.Description,
			Instruction: ac.Instruction,
			Tools: []tool.Tool{
				geminitool.GoogleSearch{},
				// agenttool.New()
			},
		})
		if err != nil {
			log.Fatalf("Failed to create agent: %v", err)
		}

		ac.agent = &currentAgent
		//fmt.Println("tool agent")
	case "root-agent":
		currentAgent, err := llmagent.New(llmagent.Config{
			Name:        ac.Name,
			Model:       model,
			Description: ac.Description,
			Instruction: ac.Instruction,
			Tools:       []tool.Tool{},
		})
		if err != nil {
			log.Fatalf("Failed to create agent: %v", err)
		}

		ac.agent = &currentAgent
		// fmt.Println("root agent")
	default:
		fmt.Println("no agent type")
		os.Exit(1)

	}

}
