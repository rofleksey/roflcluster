package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type SSHConfig struct {
	Password   string `yaml:"password"`
	KeyFile    string `yaml:"keyFile"`
	Passphrase string `yaml:"passphrase"`
}

type NodeConfig struct {
	Name      string    `yaml:"name"`
	Ip        string    `yaml:"ip"`
	Domain    string    `yaml:"domain"`
	UseDocker bool      `yaml:"useDocker"`
	Ssh       SSHConfig `yaml:"ssh"`
}

type ClusterConfig struct {
	UsePrivateRepo bool `yaml:"usePrivateRepo"`
}

type UserConfig struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Bcrypt   string `yaml:"bcrypt"`
	Secret   string `yaml:"secret"`
	Email    string `yaml:"email"`
}

type Config struct {
	User       UserConfig    `yaml:"user"`
	Cluster    ClusterConfig `yaml:"cluster"`
	MainNode   NodeConfig    `yaml:"mainNode"`
	AgentNodes []NodeConfig  `yaml:"agentNodes"`
}

func ReadConfig() (*Config, error) {
	configBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	result := Config{}

	err = yaml.Unmarshal(configBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &result, nil
}
