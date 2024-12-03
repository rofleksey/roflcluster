package util

import (
	"fmt"
	"log"
	"os"
	"roflcluster/config"
	"roflcluster/step"
	"strings"
)

func InitAgentNode(mainCfg config.NodeConfig, nodeCfg config.NodeConfig) error {
	log.Printf("Connecting to agent node %s...", nodeCfg.Name)

	client, err := OpenSSH(nodeCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	defer client.Close()

	if _, err := client.Run("bash -c 'test -f /usr/local/bin/k3s-uninstall.sh'"); err == nil {
		log.Println("Agent already exists, k3s installation skipped")

		return nil
	}

	tokenBytes, err := os.ReadFile("main-node-token")
	if err != nil {
		return fmt.Errorf("failed to read main node token file: %w", err)
	}

	token := strings.TrimSpace(string(tokenBytes))

	log.Println("Installing K3S Agent...")

	installCmd := fmt.Sprintf("curl -sfL https://get.k3s.io | K3S_URL=\"https://%s:6443\" K3S_TOKEN=\"%s\" INSTALL_K3S_EXEC=\"--node-name=%s --node-external-ip=%s\" sh -s -", mainCfg.Ip, token, nodeCfg.Name, nodeCfg.Ip)
	stdout, err := client.Run(installCmd)
	if err != nil {
		return fmt.Errorf("failed to run k3s installation: %w %s", err, string(stdout))
	}

	return nil
}

func InitMainNode(rootCfg *config.Config, scenario *step.Scenario) error {
	log.Println("Connecting to main node...")

	client, err := OpenSSH(rootCfg.MainNode)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	defer client.Close()

	log.Println("Starting scenario execution")

	for _, step := range scenario.Steps {
		log.Printf("Executing step '%s'", step.String())

		if err := step.Execute(client, rootCfg); err != nil {
			return err
		}
	}

	return nil
}
