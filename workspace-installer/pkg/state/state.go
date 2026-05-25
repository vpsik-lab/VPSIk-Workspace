package state

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vpsik/workspace-installer/pkg/detector"
)

type ServiceState struct {
	Name    string          `json:"name"`
	Status  detector.Status `json:"status"`
	Details string          `json:"details"`
}

type State struct {
	Services []ServiceState `json:"services"`
}

func Build(detection *detector.Result) *State {
	s := &State{}
	for _, svc := range detection.Services {
		s.Services = append(s.Services, ServiceState{
			Name:    svc.Name,
			Status:  svc.Status,
			Details: svc.Details,
		})
	}
	return s
}

func (s *State) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read state: %w", err)
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal state: %w", err)
	}
	return &s, nil
}
