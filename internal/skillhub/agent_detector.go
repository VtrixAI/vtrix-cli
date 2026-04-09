package skillhub

import (
	"os"
	"path/filepath"
	"strings"
)

type AgentConfig struct {
	Name            string
	DisplayName     string
	LocalSkillsDir  string
	GlobalSkillsDir string
}

func getAgentConfigs() []AgentConfig {
	home, _ := os.UserHomeDir()

	return []AgentConfig{
		{
			Name:            "cursor",
			DisplayName:     "Cursor",
			LocalSkillsDir:  ".agents/skills",
			GlobalSkillsDir: filepath.Join(home, ".cursor", "skills"),
		},
		{
			Name:            "claude-code",
			DisplayName:     "Claude Code",
			LocalSkillsDir:  ".claude/skills",
			GlobalSkillsDir: filepath.Join(home, ".claude", "skills"),
		},
		{
			Name:            "codex",
			DisplayName:     "Codex",
			LocalSkillsDir:  ".agents/skills",
			GlobalSkillsDir: filepath.Join(home, ".codex", "skills"),
		},
		{
			Name:            "cline",
			DisplayName:     "Cline",
			LocalSkillsDir:  ".agents/skills",
			GlobalSkillsDir: filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:            "continue",
			DisplayName:     "Continue",
			LocalSkillsDir:  ".continue/skills",
			GlobalSkillsDir: filepath.Join(home, ".continue", "skills"),
		},
		{
			Name:            "openclaw",
			DisplayName:     "OpenClaw",
			LocalSkillsDir:  "skills",
			GlobalSkillsDir: filepath.Join(home, ".openclaw", "skills"),
		},
	}
}

func DetectCurrentAgent() *AgentConfig {
	configs := getAgentConfigs()

	if os.Getenv("CURSOR_AGENT") != "" {
		for _, c := range configs {
			if c.Name == "cursor" {
				return &c
			}
		}
	}

	if os.Getenv("CODEX_HOME") != "" {
		for _, c := range configs {
			if c.Name == "codex" {
				return &c
			}
		}
	}

	if path := os.Getenv("PATH"); strings.Contains(path, "claude-code") {
		if os.Getenv("CURSOR_AGENT") == "" {
			for _, c := range configs {
				if c.Name == "claude-code" {
					return &c
				}
			}
		}
	}

	home, _ := os.UserHomeDir()
	agents := []struct {
		name string
		path string
	}{
		{"cursor", filepath.Join(home, ".cursor")},
		{"claude-code", filepath.Join(home, ".claude")},
		{"codex", filepath.Join(home, ".codex")},
		{"cline", filepath.Join(home, ".cline")},
		{"continue", filepath.Join(home, ".continue")},
		{"openclaw", filepath.Join(home, ".openclaw")},
	}

	for _, agent := range agents {
		if _, err := os.Stat(agent.path); err == nil {
			for _, c := range configs {
				if c.Name == agent.name {
					return &c
				}
			}
		}
	}

	return nil
}

func GetInstallDir(global bool, agent *AgentConfig) string {
	if global {
		if agent != nil {
			return agent.GlobalSkillsDir
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".agents", "skills")
	}

	localDir := ".agents/skills"
	if agent != nil {
		localDir = agent.LocalSkillsDir
	}

	cwd, _ := os.Getwd()
	return filepath.Join(cwd, localDir)
}

// DetectAllInstalledAgents returns all Agent configurations that are installed on the system
func DetectAllInstalledAgents() []AgentConfig {
	var installed []AgentConfig

	for _, config := range getAgentConfigs() {
		// Check if the Agent's root directory exists (e.g., ~/.claude, ~/.cursor)
		agentRoot := filepath.Dir(config.GlobalSkillsDir)
		if _, err := os.Stat(agentRoot); err == nil {
			installed = append(installed, config)
		}
	}

	return installed
}
