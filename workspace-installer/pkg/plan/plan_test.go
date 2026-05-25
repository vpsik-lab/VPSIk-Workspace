package plan

import (
	"testing"

	"github.com/vpsik/workspace-installer/pkg/detector"
	"github.com/vpsik/workspace-installer/pkg/state"
)

func TestBuild_AllInstalled(t *testing.T) {
	svcState := &state.State{
		Services: []state.ServiceState{
			{Name: "gitea", Status: detector.StatusInstalled},
			{Name: "ollama", Status: detector.StatusInstalled},
		},
	}

	p := Build(svcState, []string{"gitea", "ollama"})

	if p.HasChanges() {
		t.Error("expected no changes when all installed")
	}

	if len(p.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(p.Items))
	}

	for _, item := range p.Items {
		if item.Action != ActionSkip {
			t.Errorf("expected skip for %s, got %s", item.Service, item.Action)
		}
	}

	expected := "0 to install, 2 already installed"
	if s := p.Summary(); s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}

func TestBuild_AllMissing(t *testing.T) {
	svcState := &state.State{
		Services: []state.ServiceState{
			{Name: "gitea", Status: detector.StatusMissing, Details: "Not detected"},
			{Name: "ollama", Status: detector.StatusMissing, Details: "Not detected"},
		},
	}

	p := Build(svcState, []string{"gitea", "ollama"})

	if !p.HasChanges() {
		t.Error("expected changes when all missing")
	}

	if len(p.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(p.Items))
	}

	for _, item := range p.Items {
		if item.Action != ActionInstall {
			t.Errorf("expected install for %s, got %s", item.Service, item.Action)
		}
	}
}

func TestBuild_Mixed(t *testing.T) {
	svcState := &state.State{
		Services: []state.ServiceState{
			{Name: "gitea", Status: detector.StatusInstalled},
			{Name: "ollama", Status: detector.StatusMissing, Details: "Not found"},
		},
	}

	p := Build(svcState, []string{"gitea", "ollama"})

	if !p.HasChanges() {
		t.Error("expected changes")
	}

	if len(p.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(p.Items))
	}

	if p.Items[0].Action != ActionSkip || p.Items[0].Service != "gitea" {
		t.Error("expected gitea to be skipped")
	}
	if p.Items[1].Action != ActionInstall || p.Items[1].Service != "ollama" {
		t.Error("expected ollama to be installed")
	}
}

func TestBuild_NotFoundInState(t *testing.T) {
	svcState := &state.State{
		Services: []state.ServiceState{
			{Name: "gitea", Status: detector.StatusInstalled},
		},
	}

	p := Build(svcState, []string{"gitea", "ollama"})

	if len(p.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(p.Items))
	}

	if p.Items[1].Service != "ollama" || p.Items[1].Action != ActionInstall {
		t.Errorf("expected ollama install, got %s %s", p.Items[1].Service, p.Items[1].Action)
	}
	if p.Items[1].Reason != "Not detected in environment" {
		t.Errorf("expected 'Not detected in environment', got %q", p.Items[1].Reason)
	}
}

func TestBuild_Empty(t *testing.T) {
	p := Build(&state.State{}, []string{})
	if len(p.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(p.Items))
	}
	if p.HasChanges() {
		t.Error("expected no changes for empty plan")
	}
}

func TestBuild_MissingWithDetails(t *testing.T) {
	svcState := &state.State{
		Services: []state.ServiceState{
			{Name: "coolify", Status: detector.StatusMissing, Details: "Port 8000 not open"},
		},
	}

	p := Build(svcState, []string{"coolify"})
	if len(p.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(p.Items))
	}
	if p.Items[0].Reason != "Missing — Port 8000 not open" {
		t.Errorf("expected reason with details, got %q", p.Items[0].Reason)
	}
}

func TestSummary(t *testing.T) {
	tests := []struct {
		name     string
		plan     *Plan
		expected string
	}{
		{"all install", &Plan{Items: []Item{{Action: ActionInstall}, {Action: ActionInstall}}}, "2 to install, 0 already installed"},
		{"all skip", &Plan{Items: []Item{{Action: ActionSkip}, {Action: ActionSkip}}}, "0 to install, 2 already installed"},
		{"mixed", &Plan{Items: []Item{{Action: ActionInstall}, {Action: ActionSkip}}}, "1 to install, 1 already installed"},
		{"empty", &Plan{}, "0 to install, 0 already installed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if s := tt.plan.Summary(); s != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, s)
			}
		})
	}
}
