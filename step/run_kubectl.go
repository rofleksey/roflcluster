package step

import (
	"fmt"
	"github.com/melbahja/goph"
	"os"
	"sneakerdocs/config"
	"strings"
)

type RunKubectlStep struct {
	Cmd       string `yaml:"cmd"`
	Namespace string `yaml:"namespace"`
	SaveFile  string `yaml:"saveFile"`
}

func (s *RunKubectlStep) Execute(client *goph.Client, cfg *config.Config) error {
	cmd := fmt.Sprintf("kubectl --kubeconfig /etc/rancher/k3s/k3s.yaml %s --namespace %s", s.Cmd, s.Namespace)
	stdout, err := client.Run(cmd)
	if err != nil {
		if !strings.Contains(string(stdout), "already exists") {
			return fmt.Errorf("failed to create run kubectl command: %w %s", err, string(stdout))
		}
	}

	if s.SaveFile != "" {
		err = os.WriteFile(s.SaveFile, stdout, 0644)
		if err != nil {
			return fmt.Errorf("failed to save kubectl output to file: %w", err)
		}
	}

	return nil
}

func (s *RunKubectlStep) String() string {
	return fmt.Sprintf("Run kubectl %s --namespace %s", s.Cmd, s.Namespace)
}
