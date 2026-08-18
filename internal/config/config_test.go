package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	for _, k := range []string{"PORT", "ADDR", "MAX_PAGE_SIZE"} {
		os.Unsetenv(k)
	}
	cfg := Load()
	if cfg.Addr != ":8080" {
		t.Fatalf("addr = %s", cfg.Addr)
	}
	if cfg.MaxPageSize != 100 {
		t.Fatalf("max page size = %d", cfg.MaxPageSize)
	}
}

func TestLoadWithEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("MAX_PAGE_SIZE", "50")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("MAX_PAGE_SIZE")
	}()
	cfg := Load()
	if cfg.Addr != ":9090" || cfg.MaxPageSize != 50 {
		t.Fatalf("cfg = %+v", cfg)
	}
}

func TestLoadAddrOverridesPort(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("ADDR", "127.0.0.1:7070")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ADDR")
	}()
	if cfg := Load(); cfg.Addr != "127.0.0.1:7070" {
		t.Fatalf("addr = %s", cfg.Addr)
	}
}

func TestLoadInvalidMaxPageSizeFallsBack(t *testing.T) {
	os.Setenv("MAX_PAGE_SIZE", "abc")
	defer os.Unsetenv("MAX_PAGE_SIZE")
	if cfg := Load(); cfg.MaxPageSize != 100 {
		t.Fatalf("max page size = %d, want 100", cfg.MaxPageSize)
	}
}

func TestLoadNonPositiveMaxPageSizeFallsBack(t *testing.T) {
	os.Setenv("MAX_PAGE_SIZE", "-5")
	defer os.Unsetenv("MAX_PAGE_SIZE")
	if cfg := Load(); cfg.MaxPageSize != 100 {
		t.Fatalf("max page size = %d, want 100", cfg.MaxPageSize)
	}
}

func TestStringMethod(t *testing.T) {
	for _, k := range []string{"PORT", "ADDR", "MAX_PAGE_SIZE"} {
		os.Unsetenv(k)
	}
	cfg := Load()
	if cfg.String() == "" {
		t.Fatalf("String() should not be empty")
	}
}
