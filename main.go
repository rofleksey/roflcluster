package main

import (
	"flag"
	"log"
	"roflcluster/config"
	"roflcluster/step"
	"roflcluster/util"
)

func main() {
	destroyClusterFlag := flag.Bool("destroy", false, "destroy cluster")
	flag.Parse()

	rootCfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	scenario := step.CreateScenario(rootCfg)

	if *destroyClusterFlag {
		err = util.DestroyExistingCluster(rootCfg)
		if err != nil {
			log.Fatalf("failed to destroy existing cluster: %s", err.Error())
		}

		log.Println("Cluster has been destroyed")

		return
	}

	err = util.InitMainNode(rootCfg, scenario)
	if err != nil {
		log.Fatalf("failed to init main node: %s", err.Error())
	}

	for _, nodeCfg := range rootCfg.AgentNodes {
		err = util.InitAgentNode(rootCfg.MainNode, nodeCfg)
		if err != nil {
			log.Fatalf("failed to init agent node %s: %s", nodeCfg.Name, err.Error())
		}
	}

	log.Println("KubeConfig saved to 'k3s.yaml' file")
	log.Println("Dashboard token saved to 'dashboard-token' file")
	log.Println("Cluster successfully upgraded!")
}
