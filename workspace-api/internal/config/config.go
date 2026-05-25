package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type APIConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Auth     AuthConfig     `yaml:"auth"`
	Services ServicesConfig `yaml:"services"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret"`
	Users     []User `yaml:"users"`
}

type User struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
}

type ServicesConfig struct {
	Gitea    ServiceEndpoint `yaml:"gitea"`
	Coolify  ServiceEndpoint `yaml:"coolify"`
	Ollama   ServiceEndpoint `yaml:"ollama"`
	OpenCode ServiceEndpoint `yaml:"opencode"`
	Restic   ResticConfig    `yaml:"restic"`
}

type ResticConfig struct {
	Binary   string `yaml:"binary"`
	RepoURL  string `yaml:"repo_url"`
	Password string `yaml:"password"`
}

type ServiceEndpoint struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

func Load(path string) (*APIConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &APIConfig{
		Server: ServerConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func (c *APIConfig) FindUser(username string) *User {
	for _, u := range c.Auth.Users {
		if u.Username == username {
			return &u
		}
	}
	return nil
}
