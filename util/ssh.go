package util

import (
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"roflcluster/config"
)

func OpenSSH(nodeConfig config.NodeConfig) (*goph.Client, error) {
	var auth goph.Auth
	var err error

	sshCfg := nodeConfig.Ssh

	if nodeConfig.Ssh.KeyFile != "" {
		auth, err = goph.Key(sshCfg.KeyFile, sshCfg.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to read ssh key file: %w", err)
		}
	} else {
		auth = goph.Password(sshCfg.Password)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     "root",
		Addr:     nodeConfig.Ip,
		Port:     22,
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ssh client: %w", err)
	}

	return client, nil
}
