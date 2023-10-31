package util

import (
	"fmt"
	"log"
	"sneakerdocs/config"
	"strings"
)

func destroyExistingNode(nodeCfg config.NodeConfig) error {
	client, err := OpenSSH(nodeCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	defer client.Close()

	stdout, err := client.Run("bash -c '/usr/local/bin/k3s-uninstall.sh'")
	if err != nil {
		if strings.Contains(string(stdout), "No such file or directory") {
			return nil
		}

		return fmt.Errorf("failed to run k3s-uninstall: %w %s", err, stdout)
	}

	return nil
}

func DestroyExistingCluster(cfg *config.Config) error {
	log.Println("Destroying main node...")
	if err := destroyExistingNode(cfg.MainNode); err != nil {
		return fmt.Errorf("failed to destroy main node: %w", err)
	}

	for _, nodeCfg := range cfg.AgentNodes {
		log.Printf("Destroying node %s...", nodeCfg.Name)

		if err := destroyExistingNode(nodeCfg); err != nil {
			return fmt.Errorf("failed to destroy node %s: %w", nodeCfg.Name, err)
		}
	}

	return nil
}
