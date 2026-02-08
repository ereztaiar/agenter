package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/adk/tool/geminitool"

	"regexp"
	"strings"
)

type CustomFunction struct {
	Name      string `yaml:"name"`
	Path      string `yaml:"path"`
	InputArgs string `yaml:"input_args"`
	Results   string `yaml:"output"`
}
type CustomFunctions []CustomFunction
type PredefinedFunctions []string

type Tools struct {
	Agents              Agents              `yaml:"agents,omitempty"`
	CustomFunctions     CustomFunctions     `yaml:"custom_functions,omitempty"`
	PredefinedFunctions PredefinedFunctions `yaml:"predefined_functions,omitempty"`
}

type FieldDefinition struct {
	Name     string
	Type     string
	Optional bool
}

// parseInputArgs parses the input_args string and returns field definitions
func parseInputArgs(inputArgs string) ([]FieldDefinition, error) {
	// Remove outer braces and whitespace
	inputArgs = strings.TrimSpace(inputArgs)
	inputArgs = strings.Trim(inputArgs, "\"")
	inputArgs = strings.Trim(inputArgs, "{}")

	var fields []FieldDefinition

	// Split by comma
	parts := strings.Split(inputArgs, ",")

	// Regex to match "fieldName?: type" or "fieldName: type"
	fieldRegex := regexp.MustCompile(`^\s*(\w+)(\?)?:\s*(\w+)\s*$`)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		matches := fieldRegex.FindStringSubmatch(part)
		if matches == nil {
			return nil, fmt.Errorf("invalid field format: %s", part)
		}

		fieldName := matches[1]
		isOptional := matches[2] == "?"
		fieldType := matches[3]

		fields = append(fields, FieldDefinition{
			Name:     fieldName,
			Type:     fieldType,
			Optional: isOptional,
		})
	}

	return fields, nil
}

func (cf *CustomFunctions) generate() []tool.Tool {
	tools := []tool.Tool{}
	for _, function := range *cf {
		var functionTool tool.Tool

		type Input map[string]interface{}

		// Parse the input args to understand required vs optional fields
		fields, err := parseInputArgs(function.InputArgs)
		if err != nil {
			log.Printf("failed to parse input args for %s: %v", function.Name, err)
			continue
		}

		// Generate description with required/optional field info
		description := fmt.Sprintf("runs %s python method.\n", function.Name)
		description += "Parameters:\n"
		for _, field := range fields {
			optionalTag := ""
			if field.Optional {
				optionalTag = " (optional)"
			}
			description += fmt.Sprintf("- %s (%s)%s\n", field.Name, field.Type, optionalTag)
		}

		executeFnAsync := func(ctx tool.Context, input map[string]interface{}) (map[string]interface{}, error) {

			// Validate required fields
			for _, field := range fields {
				if !field.Optional {
					if _, exists := input[field.Name]; !exists {
						return nil, fmt.Errorf("required field '%s' is missing", field.Name)
					}
				}
			}

			var stringValues []string
			for _, value := range input {
				if strValue, ok := value.(string); ok {
					stringValues = append(stringValues, strValue)
				}
			}
			args := append([]string{function.Path}, stringValues...)

			// Execute Python script with JSON input
			cmd := exec.Command("python3", args...)
			output, err := cmd.Output()

			var response map[string]interface{}
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(output, &response)
			if err != nil {
				return nil, err
			}

			return response, nil
		}

		functionTool, err = functiontool.New(
			functiontool.Config{
				Name:        function.Name,
				Description: description,
			},
			executeFnAsync,
		)
		if err != nil {
			log.Fatalln("failed to create long running tool: %w", err)
		}

		tools = append(tools, functionTool)
	}
	return tools
}

func (pf *PredefinedFunctions) generate() []tool.Tool {
	tools := []tool.Tool{}
	for _, functionName := range *pf {
		var functionTool tool.Tool

		switch functionName {
		case "GoogleSearch":
			functionTool = geminitool.GoogleSearch{}
		default:

		}
		tools = append(tools, functionTool)
	}
	return tools
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

	customTools := t.CustomFunctions.generate()
	tools = append(tools, customTools...)
	predefinedTools := t.PredefinedFunctions.generate()
	tools = append(tools, predefinedTools...)

	return tools
}
