package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Config holds all application configuration loaded from .env
type Config struct {
	// App identity
	AppName     string
	AppVersion  string
	BuildNumber string
	CurrentYear string

	// Window
	WindowTitle  string
	WindowWidth  int
	WindowHeight int

	// Database
	DBDirName string
	DBFile    string
}

var loaded *Config

// Load reads .env from the project root (executable directory or workspace root).
// Subsequent calls return the cached instance.
func Load() (*Config, error) {
	if loaded != nil {
		return loaded, nil
	}

	envPath := resolveEnvPath()
	vals, err := parseEnvFile(envPath)
	if err != nil {
		// .env not found is acceptable — fall back to OS env / defaults
		vals = map[string]string{}
	}

	// Merge OS environment (OS takes precedence over .env file)
	for k, v := range vals {
		if _, exists := os.LookupEnv(k); !exists {
			os.Setenv(k, v)
		}
	}

	loaded = &Config{
		AppName:      getEnv("APP_NAME", "Personal Cockpit"),
		AppVersion:   getEnv("APP_VERSION", "1.0.0"),
		BuildNumber:  getEnv("BUILD_NUMBER", ""),
		CurrentYear:  getEnv("CURRENT_YEAR", "2026"),
		WindowTitle:  getEnv("WINDOW_TITLE", "Personal Cockpit"),
		WindowWidth:  getEnvInt("WINDOW_WIDTH", 1024),
		WindowHeight: getEnvInt("WINDOW_HEIGHT", 768),
		DBDirName:    getEnv("DB_DIR_NAME", "Personal Cockpit"),
		DBFile:       getEnv("DB_FILE", "cockpit.db"),
	}

	return loaded, nil
}

// MustLoad calls Load and panics on error.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("config: %v", err))
	}
	return cfg
}

// FullVersion returns the version string, optionally suffixed with build number.
func (c *Config) FullVersion() string {
	if c.BuildNumber != "" {
		return c.AppVersion + "b" + c.BuildNumber
	}
	return c.AppVersion
}

// resolveEnvPath finds the .env file relative to the binary or source root.
func resolveEnvPath() string {
	// 1. Beside the running executable (production)
	exe, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exe), ".env")
		if fileExists(candidate) {
			return candidate
		}
	}

	// 2. Walk up from source file location (development — __FILE__ trick via runtime)
	_, srcFile, _, ok := runtime.Caller(0)
	if ok {
		// srcFile is config/config.go → go two levels up to project root
		root := filepath.Dir(filepath.Dir(srcFile))
		candidate := filepath.Join(root, ".env")
		if fileExists(candidate) {
			return candidate
		}
	}

	// 3. Current working directory fallback
	return ".env"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// parseEnvFile reads a .env file and returns key=value pairs.
// Supports comments (#), blank lines, quoted values, and inline comments.
func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// Strip inline comment
		if i := strings.Index(val, " #"); i >= 0 {
			val = strings.TrimSpace(val[:i])
		}

		// Strip surrounding quotes
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		result[key] = val
	}
	return result, scanner.Err()
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
