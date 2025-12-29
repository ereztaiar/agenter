package agent

import (
	"google.golang.org/adk/agent"
)

type SubAgents map[string]*agent.Agent

// UnmarshalYAML implements custom unmarshaling for SubAgents from a YAML array
func (s *SubAgents) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var names []string
	if err := unmarshal(&names); err != nil {
		return err
	}

	*s = make(SubAgents, len(names))
	for _, name := range names {
		(*s)[name] = nil // agents are not yet initialzed
	}
	return nil
}

// MarshalYAML implements custom marshaling for SubAgents to a YAML array
func (s SubAgents) MarshalYAML() (interface{}, error) {
	if len(s) == 0 {
		return nil, nil
	}

	names := make([]string, 0, len(s))
	for name := range s {
		names = append(names, name)
	}

	return names, nil
}

func (s *SubAgents) buildSubagents(){
	
}
