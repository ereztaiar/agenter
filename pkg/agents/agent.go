package agent

type Agents []string
type Functions []string

// Name  []string
// agent *[]agent.Agent

type AgentType string
type DependsOn string
type Name string
type Model string
type Workflow string
type Description string
type Instruction string
type ApiKey string
type Arguments string
type OutputKey *string

// var patterns = map[string]func(query string) string{
// 	"parallel":   func(query string) string { return "searched " + query },
// 	"sequential": func(query string) string { return "helped " + query },
// 	"loop":       func(query string) string { return "helped " + query },
// }

// if cmdFunc, exists := commands[command]; exists {
//     result := cmdFunc(query)
//     // ...
// }


