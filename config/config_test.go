package config

import (
	"os"
	"testing"
)

// reset limpa o cache entre testes para evitar interferência.
func reset() {
	loaded = nil
}

func TestDefaultValues(t *testing.T) {
	reset()

	cfg := MustLoad()

	if cfg.AppName == "" {
		t.Error("AppName não pode ser vazio")
	}
	if cfg.AppVersion == "" {
		t.Error("AppVersion não pode ser vazio")
	}
	if cfg.WindowWidth <= 0 {
		t.Errorf("WindowWidth deve ser positivo, got %d", cfg.WindowWidth)
	}
	if cfg.WindowHeight <= 0 {
		t.Errorf("WindowHeight deve ser positivo, got %d", cfg.WindowHeight)
	}
	if cfg.DBDirName == "" {
		t.Error("DBDirName não pode ser vazio")
	}
	if cfg.DBFile == "" {
		t.Error("DBFile não pode ser vazio")
	}
}

func TestEnvOverride(t *testing.T) {
	reset()
	os.Setenv("APP_NAME", "Cockpit de Teste")
	os.Setenv("WINDOW_WIDTH", "1280")
	t.Cleanup(func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("WINDOW_WIDTH")
		reset()
	})

	cfg := MustLoad()

	if cfg.AppName != "Cockpit de Teste" {
		t.Errorf("esperado 'Cockpit de Teste', got '%s'", cfg.AppName)
	}
	if cfg.WindowWidth != 1280 {
		t.Errorf("esperado 1280, got %d", cfg.WindowWidth)
	}
}

func TestFullVersionWithBuildNumber(t *testing.T) {
	reset()
	os.Setenv("APP_VERSION", "2.0.0")
	os.Setenv("BUILD_NUMBER", "999")
	t.Cleanup(func() {
		os.Unsetenv("APP_VERSION")
		os.Unsetenv("BUILD_NUMBER")
		reset()
	})

	cfg := MustLoad()
	want := "2.0.0b999"
	if cfg.FullVersion() != want {
		t.Errorf("esperado '%s', got '%s'", want, cfg.FullVersion())
	}
}

func TestFullVersionWithoutBuildNumber(t *testing.T) {
	reset()
	os.Setenv("APP_VERSION", "1.5.0")
	os.Setenv("BUILD_NUMBER", "")
	t.Cleanup(func() {
		os.Unsetenv("APP_VERSION")
		os.Unsetenv("BUILD_NUMBER")
		reset()
	})

	cfg := MustLoad()
	want := "1.5.0"
	if cfg.FullVersion() != want {
		t.Errorf("esperado '%s', got '%s'", want, cfg.FullVersion())
	}
}

func TestCacheReturnsSameInstance(t *testing.T) {
	reset()
	a := MustLoad()
	b := MustLoad()
	if a != b {
		t.Error("MustLoad deve retornar a mesma instância em chamadas consecutivas")
	}
}

func TestInvalidWindowWidthFallsBackToDefault(t *testing.T) {
	reset()
	os.Setenv("WINDOW_WIDTH", "nao-e-numero")
	t.Cleanup(func() {
		os.Unsetenv("WINDOW_WIDTH")
		reset()
	})

	cfg := MustLoad()
	if cfg.WindowWidth != 1024 {
		t.Errorf("esperado fallback 1024, got %d", cfg.WindowWidth)
	}
}
