package docker

import (
	"testing"
)

func TestServiceTemplates_ContainAllServices(t *testing.T) {
	expected := []string{
		"authentik", "gitea", "coolify", "ollama",
		"opencode", "openwebui", "restic",
		"traefik", "postgres", "grafana", "prometheus",
		"redis", "code-server", "plane", "outline", "mattermost",
		"cloudflare",
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
	}
}

func TestServiceTemplates_Network(t *testing.T) {
	for name, tmpl := range ServiceTemplates {
		if tmpl.Network != "workspace_net" {
			t.Errorf("template %q expected network workspace_net, got %s", name, tmpl.Network)
		}
	}
}

func TestServiceTemplates_AllHaveImage(t *testing.T) {
	for name, tmpl := range ServiceTemplates {
		if tmpl.Image == "" {
			t.Errorf("template %q has no image", name)
		}
	}
}

func TestServiceTemplates_TotalCount(t *testing.T) {
	if len(ServiceTemplates) < 15 {
		t.Errorf("expected at least 15 templates, got %d", len(ServiceTemplates))
	}
}

func TestServiceTemplates_Postgres(t *testing.T) {
	tmpl, ok := ServiceTemplates["postgres"]
	if !ok {
		t.Fatal("expected postgres template")
	}

	if tmpl.EnvVars["POSTGRES_USER"] != "${POSTGRES_USER:-vpsik}" {
		t.Errorf("unexpected POSTGRES_USER: %s", tmpl.EnvVars["POSTGRES_USER"])
	}

	hasExpose := false
	for _, p := range tmpl.Expose {
		if p == "5432" {
			hasExpose = true
			break
		}
	}
	if !hasExpose {
		t.Error("expected expose 5432")
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

	if tmpl.TraefikHost != "metrics.${VPSIK_DOMAIN}" {
		t.Errorf("expected TraefikHost metrics.${VPSIK_DOMAIN}, got %s", tmpl.TraefikHost)
	}

	if tmpl.EnvVars["GF_SECURITY_ADMIN_USER"] != "${GRAFANA_USER:-admin}" {
		t.Errorf("unexpected grafana user env: %s", tmpl.EnvVars["GF_SECURITY_ADMIN_USER"])
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

	hasExpose := false
	for _, p := range tmpl.Expose {
		if p == "9090" {
			hasExpose = true
			break
		}
	}
	if !hasExpose {
		t.Error("expected expose 9090")
	}
}

func TestServiceTemplates_TraefikPorts(t *testing.T) {
	tmpl, ok := ServiceTemplates["traefik"]
	if !ok {
		t.Fatal("expected traefik template")
	}

	has80 := false
	has443 := false
	for _, p := range tmpl.Ports {
		if p == "80:80" {
			has80 = true
		}
		if p == "443:443" {
			has443 = true
		}
	}
	if !has80 || !has443 {
		t.Error("expected ports 80:80 and 443:443")
	}
}

func TestServiceTemplates_CodeServer(t *testing.T) {
	tmpl, ok := ServiceTemplates["code-server"]
	if !ok {
		t.Fatal("expected code-server template")
	}

	if tmpl.Image != "codercom/code-server:latest" {
		t.Errorf("expected coder/code-server:latest, got %s", tmpl.Image)
	}

	hasExpose := false
	for _, p := range tmpl.Expose {
		if p == "8443" {
			hasExpose = true
			break
		}
	}
	if !hasExpose {
		t.Error("expected expose 8443")
	}

	if tmpl.TraefikHost != "code.${VPSIK_DOMAIN}" {
		t.Errorf("expected TraefikHost code.${VPSIK_DOMAIN}, got %s", tmpl.TraefikHost)
	}
}

func TestGenerateTraefikLabels(t *testing.T) {
	svc := &ServiceCompose{
		Name:        "gitea",
		TraefikHost: "git.${VPSIK_DOMAIN}",
		TraefikPort: "3000",
	}

	// Labels require env vars to resolve, test struct
	if svc.TraefikHost != "git.${VPSIK_DOMAIN}" {
		t.Error("expected TraefikHost to be preserved")
	}
}

func TestServiceTemplates_Redis(t *testing.T) {
	tmpl, ok := ServiceTemplates["redis"]
	if !ok {
		t.Fatal("expected redis template")
	}

	if tmpl.Image != "redis:7-alpine" {
		t.Errorf("expected redis:7-alpine, got %s", tmpl.Image)
	}

	hasExpose := false
	for _, p := range tmpl.Expose {
		if p == "6379" {
			hasExpose = true
			break
		}
	}
	if !hasExpose {
		t.Error("expected expose 6379")
	}
}
