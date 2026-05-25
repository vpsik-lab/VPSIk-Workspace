package state

import (
	"github.com/vpsik/workspace-installer/pkg/detector"
)

type ServiceState struct {
	Name    string         `json:"name"`
	Status  detector.Status `json:"status"`
	Details string         `json:"details"`
}

type State struct {
	Services []ServiceState `json:"services"`
}

func Build(detection *detector.Result) *State {
	state := &State{}
	for _, svc := range detection.Services {
		state.Services = append(state.Services, ServiceState{
			Name:    svc.Name,
			Status:  svc.Status,
			Details: svc.Details,
		})
	}
	return state
}
