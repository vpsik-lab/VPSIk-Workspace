package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func TestCodeServerClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("expected /healthz, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewCodeServerClient(srv.URL)
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCodeServerClient_CheckHealth_Error(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer srv.Close()

	c := NewCodeServerClient(srv.URL)
	err := c.CheckHealth(context.Background())
	if err == nil {
		t.Error("expected error for 500")
	}
}

func TestCodeServerClient_New_EmptyURL(t *testing.T) {
	c := NewCodeServerClient("")
	if c.baseURL != "" {
		t.Errorf("expected empty baseURL, got %s", c.baseURL)
	}
	err := c.CheckHealth(context.Background())
	if err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestPlaneClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/health" {
			t.Errorf("expected /api/v1/health, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewPlaneClient(srv.URL)
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestPlaneClient_CheckHealth_BadStatus(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer srv.Close()

	c := NewPlaneClient(srv.URL)
	err := c.CheckHealth(context.Background())
	if err == nil {
		t.Error("expected error for 404")
	}
}

func TestOutlineClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/health" {
			t.Errorf("expected /api/health, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewOutlineClient(srv.URL, "token123")
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestOutlineClient_Token(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewOutlineClient(srv.URL, "test-token")
	if c.token != "test-token" {
		t.Errorf("expected test-token, got %s", c.token)
	}
}

func TestMattermostClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v4/system/health" {
			t.Errorf("expected /api/v4/system/health, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewMattermostClient(srv.URL, "token")
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestMattermostClient_CheckHealth_Error(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	defer srv.Close()

	c := NewMattermostClient(srv.URL, "")
	err := c.CheckHealth(context.Background())
	if err == nil {
		t.Error("expected error for 503")
	}
}

func TestOpenCodeClient_New_NilWhenEmptyURL(t *testing.T) {
	c := NewOpenCodeClient("", "")
	if c != nil {
		t.Error("expected nil when baseURL is empty")
	}
}

func TestOpenCodeClient_New_NotNil(t *testing.T) {
	c := NewOpenCodeClient("http://opencode:8080", "key")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.apiKey != "key" {
		t.Errorf("expected key, got %s", c.apiKey)
	}
}

func TestOpenCodeClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/health" {
			t.Errorf("expected /api/health, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewOpenCodeClient(srv.URL, "key")
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestOpenCodeClient_Chat_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("expected /api/chat, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Errorf("expected Bearer testkey, got %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"reply":"Hello!"}`))
	})
	defer srv.Close()

	c := NewOpenCodeClient(srv.URL, "testkey")
	reply, err := c.Chat(context.Background(), "Hi", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "Hello!" {
		t.Errorf("expected Hello!, got %s", reply)
	}
}

func TestOpenCodeClient_Chat_Error(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`server error`))
	})
	defer srv.Close()

	c := NewOpenCodeClient(srv.URL, "key")
	_, err := c.Chat(context.Background(), "Hi", "", "")
	if err == nil {
		t.Error("expected error for 500")
	}
}

func TestGiteaClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/healthz" {
			t.Errorf("expected /api/healthz, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewGiteaClient(srv.URL, "token")
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestGiteaClient_CheckHealth_Error(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer srv.Close()

	c := NewGiteaClient(srv.URL, "")
	err := c.CheckHealth(context.Background())
	if err == nil {
		t.Error("expected error for 500")
	}
}

func TestCoolifyClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/health" {
			t.Errorf("expected /api/v1/health, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Errorf("expected Bearer testtoken, got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewCoolifyClient(srv.URL, "testtoken")
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestOllamaClient_CheckHealth_OK(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("expected /api/tags, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"models":[]}`))
	})
	defer srv.Close()

	c := NewOllamaClient(srv.URL)
	err := c.CheckHealth(context.Background())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestOllamaClient_ListModels(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("expected /api/tags, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"models":[{"name":"llama3.2:latest","size":100,"digest":"abc","modified_at":"2024-01-01T00:00:00Z"}]}`))
	})
	defer srv.Close()

	c := NewOllamaClient(srv.URL)
	models, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}
	if models[0].Name != "llama3.2:latest" {
		t.Errorf("expected llama3.2:latest, got %s", models[0].Name)
	}
}

func TestOllamaClient_Chat(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("expected /api/chat, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"model":"llama3.2","message":{"role":"assistant","content":"Hello!"},"done":true}`))
	})
	defer srv.Close()

	c := NewOllamaClient(srv.URL)
	reply, err := c.Chat(context.Background(), "llama3.2", []ChatMessage{{Role: "user", Content: "Hi"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "Hello!" {
		t.Errorf("expected Hello!, got %s", reply)
	}
}

func TestOllamaClient_Chat_Error(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer srv.Close()

	c := NewOllamaClient(srv.URL)
	_, err := c.Chat(context.Background(), "model", nil)
	if err == nil {
		t.Error("expected error for 500")
	}
}

func TestOllamaClient_PullModel(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pull" {
			t.Errorf("expected /api/pull, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewOllamaClient(srv.URL)
	err := c.PullModel(context.Background(), "llama3.2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestOllamaClient_DeleteModel(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/delete" {
			t.Errorf("expected /api/delete, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewOllamaClient(srv.URL)
	err := c.DeleteModel(context.Background(), "model")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCoolifyClient_ListProjects(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/projects" {
			t.Errorf("expected /api/v1/projects, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":"1","uuid":"abc","name":"Test","description":"desc"}]`))
	})
	defer srv.Close()

	c := NewCoolifyClient(srv.URL, "token")
	projects, err := c.ListProjects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != "Test" {
		t.Errorf("expected Test, got %s", projects[0].Name)
	}
}

func TestCoolifyClient_ListProjects_Error(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	defer srv.Close()

	c := NewCoolifyClient(srv.URL, "bad-token")
	_, err := c.ListProjects(context.Background())
	if err == nil {
		t.Error("expected error for 403")
	}
}

func TestGiteaClient_ListRepos(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/user/repos" {
			t.Errorf("expected /api/v1/user/repos, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "token testtoken" {
			t.Errorf("expected 'token testtoken', got %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"name":"my-repo","full_name":"user/my-repo","private":false}]`))
	})
	defer srv.Close()

	c := NewGiteaClient(srv.URL, "testtoken")
	repos, err := c.ListRepos(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if repos[0].Name != "my-repo" {
		t.Errorf("expected my-repo, got %s", repos[0].Name)
	}
}

func TestGiteaClient_CreateWebhook(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":1,"url":"https://hook.example.com","active":true,"events":["push"],"type":"gitea"}`))
	})
	defer srv.Close()

	c := NewGiteaClient(srv.URL, "token")
	hook, err := c.CreateWebhook(context.Background(), "user/repo", "https://hook.example.com", "secret", []string{"push"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hook.ID != 1 {
		t.Errorf("expected id 1, got %d", hook.ID)
	}
	if !hook.Active {
		t.Error("expected active webhook")
	}
}

func TestNewResticClient(t *testing.T) {
	c := NewResticClient("", "/repo", "pass")
	if c.binaryPath != "restic" {
		t.Errorf("expected restic, got %s", c.binaryPath)
	}
	if c.repoURL != "/repo" {
		t.Errorf("expected /repo, got %s", c.repoURL)
	}
	if c.password != "pass" {
		t.Errorf("expected pass, got %s", c.password)
	}
}

func TestResticClient_Env(t *testing.T) {
	c := NewResticClient("restic", "rest:http://repo", "secret123")
	env := c.env()
	foundRepo := false
	foundPass := false
	for _, e := range env {
		if e == "RESTIC_REPOSITORY=rest:http://repo" {
			foundRepo = true
		}
		if e == "RESTIC_PASSWORD=secret123" {
			foundPass = true
		}
	}
	if !foundRepo {
		t.Error("expected RESTIC_REPOSITORY in env")
	}
	if !foundPass {
		t.Error("expected RESTIC_PASSWORD in env")
	}
}

func TestGiteaClient_ListIssues(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/user/issues" {
			t.Errorf("expected /api/v1/user/issues, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"number":1,"title":"Bug found","state":"open"}]`))
	})
	defer srv.Close()

	c := NewGiteaClient(srv.URL, "token")
	issues, err := c.ListIssues(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Title != "Bug found" {
		t.Errorf("expected 'Bug found', got %s", issues[0].Title)
	}
}

func TestCoolifyClient_Deploy(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	c := NewCoolifyClient(srv.URL, "token")
	err := c.Deploy(context.Background(), "resource-123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestOllamaClient_ChatStream_Error(t *testing.T) {
	srv := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer srv.Close()

	c := NewOllamaClient(srv.URL)
	_, err := c.ChatStream(context.Background(), "model", nil)
	if err == nil {
		t.Error("expected error for 500")
	}
}
