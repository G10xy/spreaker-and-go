package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// resetViper clears viper state between tests so they don't interfere.
func resetViper() {
	viper.Reset()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Token != "" {
		t.Error("default Token should be empty")
	}
	if cfg.UserID != 0 {
		t.Error("default UserID should be 0")
	}
	if cfg.OutputFormat != "table" {
		t.Errorf("default OutputFormat = %q, want %q", cfg.OutputFormat, "table")
	}
	if cfg.APIURL != "https://api.spreaker.com" {
		t.Errorf("default APIURL = %q", cfg.APIURL)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	t.Setenv("SPREAKER_CONFIG_DIR", tmpDir)

	original := &Config{
		Token:         "test-token-123",
		UserID:        42,
		DefaultShowID: 99,
		OutputFormat:  "json",
		APIURL:        "https://custom.api.com",
	}

	if err := Save(original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	resetViper()
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Token != original.Token {
		t.Errorf("Token = %q, want %q", loaded.Token, original.Token)
	}
	if loaded.UserID != original.UserID {
		t.Errorf("UserID = %d, want %d", loaded.UserID, original.UserID)
	}
	if loaded.DefaultShowID != original.DefaultShowID {
		t.Errorf("DefaultShowID = %d, want %d", loaded.DefaultShowID, original.DefaultShowID)
	}
	if loaded.OutputFormat != original.OutputFormat {
		t.Errorf("OutputFormat = %q, want %q", loaded.OutputFormat, original.OutputFormat)
	}
	if loaded.APIURL != original.APIURL {
		t.Errorf("APIURL = %q, want %q", loaded.APIURL, original.APIURL)
	}
}

func TestSaveToken_PreservesOtherFields(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	t.Setenv("SPREAKER_CONFIG_DIR", tmpDir)

	initial := &Config{
		Token:         "old-token",
		UserID:        1,
		DefaultShowID: 55,
		OutputFormat:  "plain",
		APIURL:        "https://api.spreaker.com",
	}
	if err := Save(initial); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	resetViper()
	if err := SaveToken("new-token", 99); err != nil {
		t.Fatalf("SaveToken failed: %v", err)
	}

	resetViper()
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Token != "new-token" {
		t.Errorf("Token = %q, want %q", loaded.Token, "new-token")
	}
	if loaded.UserID != 99 {
		t.Errorf("UserID = %d, want 99", loaded.UserID)
	}
	if loaded.DefaultShowID != 55 {
		t.Errorf("DefaultShowID = %d, want 55 (should be preserved)", loaded.DefaultShowID)
	}
}

func TestGetUserID_NoUserID(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	t.Setenv("SPREAKER_CONFIG_DIR", tmpDir)

	// Save config with no user ID
	if err := Save(&Config{OutputFormat: "table", APIURL: "https://api.spreaker.com"}); err != nil {
		t.Fatal(err)
	}

	resetViper()
	_, err := GetUserID()
	if err == nil {
		t.Fatal("expected error when no user ID cached")
	}
}

func TestGetToken_NoToken(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	t.Setenv("SPREAKER_CONFIG_DIR", tmpDir)

	if err := Save(&Config{OutputFormat: "table", APIURL: "https://api.spreaker.com"}); err != nil {
		t.Fatal(err)
	}

	resetViper()
	_, err := GetToken()
	if err == nil {
		t.Fatal("expected error when no token set")
	}
}

func TestConfigFilePath_ReturnsPath(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	t.Setenv("SPREAKER_CONFIG_DIR", tmpDir)

	path := ConfigFilePath()
	if path == "(unknown)" {
		t.Error("ConfigFilePath returned (unknown)")
	}
	if path == "" {
		t.Error("ConfigFilePath returned empty string")
	}
}

func TestConfigFilePermissions(t *testing.T) {
	resetViper()
	tmpDir := t.TempDir()
	t.Setenv("SPREAKER_CONFIG_DIR", tmpDir)

	if err := Save(DefaultConfig()); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(tmpDir, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("config file permissions = %o, want 0600", perm)
	}
}

func TestConfigDir_RelativePath_Error(t *testing.T) {
	resetViper()
	t.Setenv("SPREAKER_CONFIG_DIR", "relative/path")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for relative SPREAKER_CONFIG_DIR")
	}
}
