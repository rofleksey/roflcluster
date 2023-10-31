package step

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/melbahja/goph"
	"log"
	"sneakerdocs/config"
	"text/template"
)

type InstallK3SStep struct {
}

func (s *InstallK3SStep) Execute(client *goph.Client, cfg *config.Config) error {
	if _, err := client.Run("bash -c 'test -f /usr/local/bin/k3s-uninstall.sh'"); err == nil {
		log.Println("Server already exists, k3s installation skipped")

		return nil
	}

	tmplStr := "curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=\"--write-kubeconfig-mode=644 --node-name={{.Name}} --tls-san={{.Ip}} --node-external-ip={{.Ip}} --flannel-backend=wireguard-native --flannel-external-ip\" sh -s -"

	tmpl, err := template.New("install-k3s").Funcs(sprig.FuncMap()).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to create installation cmd template: %w", err)
	}

	var renderBuffer bytes.Buffer

	err = tmpl.Execute(&renderBuffer, cfg.MainNode)
	if err != nil {
		return fmt.Errorf("failed to render installation cmd template: %w", err)
	}

	installCmd := renderBuffer.String()

	log.Println("Running " + installCmd)

	stdout, err := client.Run(installCmd)
	if err != nil {
		return fmt.Errorf("failed to run k3s installation: %w %s", err, string(stdout))
	}

	return nil
}

func (s *InstallK3SStep) String() string {
	return "Install K3S Server"
}
