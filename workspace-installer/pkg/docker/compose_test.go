package docker

import (
	"testing"
)

func TestServiceTemplates_ContainAllServices(t *testing.T) {
	expected := []string{
		"authentik", "gitea", "coolify", "ollama",
		"opencode", "openwebui", "restic",
		"traefik", "postgres", "grafana", "prometheus",
	}

	for _, name := range expected {
		tmpl, ok := ServiceTemplates[name]
		if !ok {
			t.Errorf("expected template %q to exist", name)
			continue
		}
		if tmpl.Name == "" {
			t.Errorf("template %q has empty name", name)
		}
		if tmpl.Image == "" {
			t.Errorf("template %q has empty image", name)
		}
		if tmpl.Network != "vpsik" {
			t.Errorf("template %q expected network vpsik, got %s", name, tmpl.Network)
		}
	}
}

func TestServiceTemplates_Grafana(t *testing.T) {
	tmpl, ok := ServiceTemplates["grafana"]
	if !ok {
		t.Fatal("expected grafana template")
	}

	if tmpl.Image != "grafana/grafana:latest" {
		t.Errorf("expected grafana/grafana:latest, got %s", tmpl.Image)
	}

	hasPort := false
	for _, p := range tmpl.Ports {
		if p == "3002:3000" {
			hasPort = true
			break
		}
	}
	if !hasPort {
		t.Error("expected port mapping 3002:3000")
	}

	hasVolume := false
	for _, v := range tmpl.Volumes {
		if v == "grafana-data:/var/lib/grafana" {
			hasVolume = true
			break
		}
	}
	if !hasVolume {
		t.Error("expected grafana-data volume")
	}

	if tmpl.EnvVars["GF_SECURITY_ADMIN_USER"] != "${GRAFANA_USER:-admin}" {
		t.Errorf("unexpected grafana user env: %s", tmpl.EnvVars["GF_SECURITY_ADMIN_USER"])
	}
	if tmpl.EnvVars["GF_SECURITY_ADMIN_PASSWORD"] != "${GRAFANA_PASSWORD:-admin}" {
		t.Errorf("unexpected grafana password env: %s", tmpl.EnvVars["GF_SECURITY_ADMIN_PASSWORD"])
	}
}

func TestServiceTemplates_Prometheus(t *testing.T) {
	tmpl, ok := ServiceTemplates["prometheus"]
	if !ok {
		t.Fatal("expected prometheus template")
	}

	if tmpl.Image != "prom/prometheus:latest" {
		t.Errorf("expected prom/prometheus:latest, got %s", tmpl.Image)
	}

	hasPort := false
	for _, p := range tmpl.Ports {
		if p == "9090:9090" {
			hasPort = true
			break
		}
	}
	if !hasPort {
		t.Error("expected port mapping 9090:9090")
	}

	hasVolume := false
	for _, v := range tmpl.Volumes {
		if v == "prometheus-data:/prometheus" {
			hasVolume = true
			break
		}
	}
	if !hasVolume {
		t.Error("expected prometheus-data volume")
	}
}

func TestServiceTemplates_Restic(t *testing.T) {
	tmpl, ok := ServiceTemplates["restic"]
	if !ok {
		t.Fatal("expected restic template")
	}

	if tmpl.Image != "restic/restic:latest" {
		t.Errorf("expected restic/restic:latest, got %s", tmpl.Image)
	}

	if tmpl.Network != "vpsik" {
		t.Errorf("expected network vpsik, got %s", tmpl.Network)
	}
}

func TestServiceTemplates_Postgres(t *testing.T) {
	tmpl, ok := ServiceTemplates["postgres"]
	if !ok {
		t.Fatal("expected postgres template")
	}

	if tmpl.Image != "postgres:16-alpine" {
		t.Errorf("expected postgres:16-alpine, got %s", tmpl.Image)
	}

	if tmpl.EnvVars["POSTGRES_USER"] != "${POSTGRES_USER:-vpsik}" {
		t.Errorf("unexpected POSTGRES_USER: %s", tmpl.EnvVars["POSTGRES_USER"])
	}
}

func TestServiceTemplates_AllHaveImage(t *testing.T) {
	for name, tmpl := range ServiceTemplates {
		if tmpl.Image == "" {
			t.Errorf("template %q has no image", name)
		}
	}
}

func TestServiceTemplates_AllHaveNetwork(t *testing.T) {
	for name, tmpl := range ServiceTemplates {
		if tmpl.Network != "vpsik" {
			t.Errorf("template %q has wrong network: %s", name, tmpl.Network)
		}
	}
}

func TestServiceTemplates_TotalCount(t *testing.T) {
	if len(ServiceTemplates) < 10 {
		t.Errorf("expected at least 10 templates, got %d", len(ServiceTemplates))
	}
}
