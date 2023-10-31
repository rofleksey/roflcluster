package step

import (
	"fmt"
	"github.com/melbahja/goph"
	"path/filepath"
	"sneakerdocs/config"
	"strings"
)

type ApplyHelmStep struct {
	RepoUrl        string `yaml:"repoUrl"`
	RepoName       string `yaml:"repoName"`
	ReleaseName    string `yaml:"releaseName"`
	Chart          string `yaml:"chart"`
	Namespace      string `yaml:"namespace"`
	ValuesTemplate string `yaml:"valuesTemplate"`
}

func (s *ApplyHelmStep) Execute(client *goph.Client, cfg *config.Config) error {
	cmd := fmt.Sprintf("helm --kubeconfig /etc/rancher/k3s/k3s.yaml repo add %s %s", s.RepoName, s.RepoUrl)
	stdout, err := client.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to add helm repository: %w %s", err, string(stdout))
	}

	data := TemplateData{
		Node: cfg.MainNode,
		User: cfg.User,
	}

	var valuesStr string

	if s.ValuesTemplate != "" {
		err = FormatUploadTemplate(client, s.ValuesTemplate, &data)
		if err != nil {
			return fmt.Errorf("failed to format and upload helm values: %w %s", err, string(stdout))
		}

		valuesStr = "-f " + strings.TrimSuffix(filepath.Base(s.ValuesTemplate), ".tmpl")
	}

	var namespaceStr string

	if s.Namespace == "" {
		namespaceStr = "--namespace default"
	} else {
		namespaceStr = "--create-namespace --namespace " + s.Namespace
	}

	cmd = fmt.Sprintf("helm --kubeconfig /etc/rancher/k3s/k3s.yaml upgrade -i %s %s %s %s", s.ReleaseName, s.Chart, namespaceStr, valuesStr)
	stdout, err = client.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to upgrade helm chart: %w %s", err, string(stdout))
	}

	return nil
}

func (s *ApplyHelmStep) String() string {
	return "Apply helm chart " + s.Chart
}
