package step

import (
	"fmt"
	"github.com/melbahja/goph"
	"os"
	"sneakerdocs/config"
	"strings"
)

type DownloadK3SConfigStep struct {
}

func (s *DownloadK3SConfigStep) Execute(client *goph.Client, cfg *config.Config) error {
	err := client.Download("/etc/rancher/k3s/k3s.yaml", "k3s.yaml")
	if err != nil {
		return fmt.Errorf("failed to download kubeconfig file: %w", err)
	}

	kubeBytes, err := os.ReadFile("k3s.yaml")
	if err != nil {
		return fmt.Errorf("failed to read downloaded kubeconfig file: %w", err)
	}

	replacedUrlKubeConfigStr := strings.Replace(string(kubeBytes), "https://127.0.0.1:6443", fmt.Sprintf("https://%s:6443", cfg.MainNode.Ip), 1)

	err = os.WriteFile("k3s.yaml", []byte(replacedUrlKubeConfigStr), 0644)
	if err != nil {
		return fmt.Errorf("failed to replace kubeconfig file ip: %w", err)
	}

	err = client.Download("/var/lib/rancher/k3s/server/node-token", "main-node-token")
	if err != nil {
		return fmt.Errorf("failed to download node token file: %w", err)
	}

	return nil
}

func (s *DownloadK3SConfigStep) String() string {
	return "Download K3S config"
}
