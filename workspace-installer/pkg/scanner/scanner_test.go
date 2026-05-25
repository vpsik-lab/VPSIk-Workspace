package scanner

import (
	"testing"
)

func TestDetectCoolify_ByContainer(t *testing.T) {
	containers := []ContainerInfo{
		{Name: "coolify-server", Image: "coollabsio/coolify:latest", Status: "running", Ports: []string{"0.0.0.0:8000->8000/tcp"}},
	}
	result := detectCoolify(containers, []int{8000, 80, 443})
	if result == nil {
		t.Fatal("expected coolify detection")
	}
	if !result.Installed {
		t.Error("expected installed")
	}
	if !result.Running {
		t.Error("expected running")
	}
	if result.Container != "coolify-server" {
		t.Errorf("expected coolify-server, got %s", result.Container)
	}
	if result.Port != 8000 {
		t.Errorf("expected port 8000, got %d", result.Port)
	}
	if !result.HasProxy {
		t.Error("expected HasProxy true when 80/443 in ports")
	}
}

func TestDetectCoolify_ByPort(t *testing.T) {
	result := detectCoolify(nil, []int{3000})
	if result == nil {
		t.Fatal("expected coolify detection")
	}
	if !result.Installed {
		t.Error("expected installed")
	}
	if result.Port != 3000 {
		t.Errorf("expected port 3000, got %d", result.Port)
	}
}

func TestDetectCoolify_NotFound(t *testing.T) {
	result := detectCoolify(nil, []int{22, 5432})
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}

func TestCollectErrors_NoErrors(t *testing.T) {
	result := &ScanResult{}
	result.System.RAMMB = 4096
	result.System.DiskFreeMB = 20480
	collectErrors(result)
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(result.Errors))
	}
}

func TestCollectErrors_LowMemory(t *testing.T) {
	result := &ScanResult{}
	result.System.RAMMB = 1024
	collectErrors(result)
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0] != "low memory: 1024 MB (minimum 2048 MB)" {
		t.Errorf("unexpected error: %s", result.Errors[0])
	}
}

func TestCollectErrors_LowDisk(t *testing.T) {
	result := &ScanResult{}
	result.System.RAMMB = 4096
	result.System.DiskFreeMB = 512
	collectErrors(result)
	found := false
	for _, e := range result.Errors {
		if e == "low disk space: 512 MB (minimum 10240 MB)" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected low disk space error")
	}
}

func TestCollectErrors_BothLow(t *testing.T) {
	result := &ScanResult{}
	result.System.RAMMB = 1024
	result.System.DiskFreeMB = 512
	collectErrors(result)
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestCollectErrors_ZeroValues(t *testing.T) {
	result := &ScanResult{}
	collectErrors(result)
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors when values are 0, got %d", len(result.Errors))
	}
}

func TestDetectSystem_SetsOS(t *testing.T) {
	info := &SystemInfo{}
	detectSystem(info)
	if info.OS == "" {
		t.Error("expected OS to be set")
	}
	if info.Architecture == "" {
		t.Error("expected Architecture to be set")
	}
	if info.CPU == 0 {
		t.Error("expected CPU > 0")
	}
}

func TestIsPortUsed_ZeroPort(t *testing.T) {
	if isPortUsed(0) {
		t.Error("expected port 0 to not be in use")
	}
}

func TestIsPortUsed_HighPort(t *testing.T) {
	if isPortUsed(65535) {
		t.Log("port 65535 may be in use on this system")
	}
}
