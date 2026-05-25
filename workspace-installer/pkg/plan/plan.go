package plan

import (
	"fmt"

	"github.com/vpsik/workspace-installer/pkg/state"
)

type Action string

const (
	ActionInstall Action = "install"
	ActionSkip    Action = "skip"
)

type Item struct {
	Service string `json:"service"`
	Action  Action `json:"action"`
	Reason  string `json:"reason"`
}

type Plan struct {
	Items []Item `json:"items"`
}

func Build(svcState *state.State, configServices []string) *Plan {
	plan := &Plan{}

	for _, name := range configServices {
		found := false
		for _, s := range svcState.Services {
			if s.Name == name {
				found = true
				if s.Status.String() == "installed" {
					plan.Items = append(plan.Items, Item{
						Service: name,
						Action:  ActionSkip,
						Reason:  "Already installed",
					})
				} else {
					plan.Items = append(plan.Items, Item{
						Service: name,
						Action:  ActionInstall,
						Reason:  fmt.Sprintf("Missing — %s", s.Details),
					})
				}
				break
			}
		}
		if !found {
			plan.Items = append(plan.Items, Item{
				Service: name,
				Action:  ActionInstall,
				Reason:  "Not detected in environment",
			})
		}
	}

	return plan
}

func (p *Plan) HasChanges() bool {
	for _, item := range p.Items {
		if item.Action == ActionInstall {
			return true
		}
	}
	return false
}

func (p *Plan) Summary() string {
	installCount := 0
	skipCount := 0
	for _, item := range p.Items {
		switch item.Action {
		case ActionInstall:
			installCount++
		case ActionSkip:
			skipCount++
		}
	}
	return fmt.Sprintf("%d to install, %d already installed", installCount, skipCount)
}
