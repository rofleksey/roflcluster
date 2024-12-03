package step

import (
	"github.com/melbahja/goph"
	"roflcluster/config"
)

type ScenarioStep interface {
	Execute(client *goph.Client, cfg *config.Config) error
	String() string
}

type Scenario struct {
	Steps []ScenarioStep `yaml:"steps"`
}

func CreateScenario(cfg *config.Config) *Scenario {
	steps := []ScenarioStep{
		// install k3s (if not already)
		&InstallK3SStep{},

		// download k3s access files
		&DownloadK3SConfigStep{},

		// install TLS secret
		&TLSSecretStep{},

		// install main ingress
		&ApplyTemplateStep{
			File: "templates/main-ingress.yaml.tmpl",
		},

		// install echo server
		&ApplyHelmStep{
			RepoUrl:        "https://ealenn.github.io/charts",
			RepoName:       "ealenn",
			ReleaseName:    "echoserver",
			Chart:          "ealenn/echo-server",
			Namespace:      "echoserver",
			ValuesTemplate: "templates/echo-server-values.yaml.tmpl",
		},

		// wait for echo server availability
		&HealthCheckStep{
			Url: "https://echo." + cfg.MainNode.Domain,
		},

		// install k8s dashboard
		&ApplyHelmStep{
			RepoUrl:        "https://kubernetes.github.io/dashboard/",
			RepoName:       "kubernetes-dashboard",
			ReleaseName:    "kubernetes-dashboard",
			Chart:          "kubernetes-dashboard/kubernetes-dashboard",
			Namespace:      "kubernetes-dashboard",
			ValuesTemplate: "templates/dashboard-values.yaml.tmpl",
		},

		// create k8s dashboard service account
		&RunKubectlStep{
			Cmd:       "create serviceaccount k8s-dashboard-admin",
			Namespace: "kubernetes-dashboard",
		},

		// create k8s dashboard role binding
		&ApplyTemplateStep{
			File: "templates/dashboard-admin-role-binding.yaml.tmpl",
		},

		// create 10years access token for k8s dashboard
		&RunKubectlStep{
			Cmd:       "create token k8s-dashboard-admin --duration 315360000s",
			Namespace: "kubernetes-dashboard",
			SaveFile:  "dashboard-token",
		},
	}

	if cfg.Cluster.UsePrivateRepo {
		steps = append(steps,
			// install gitea
			&ApplyHelmStep{
				RepoUrl:        "https://dl.gitea.io/charts/",
				RepoName:       "gitea-charts",
				ReleaseName:    "gitea",
				Chart:          "gitea-charts/gitea",
				Namespace:      "gitea",
				ValuesTemplate: "templates/gitea-values.yaml.tmpl",
			},
		)
	}

	steps = append(steps,
		// install argocd
		&ApplyHelmStep{
			RepoUrl:        "https://argoproj.github.io/argo-helm",
			RepoName:       "argo",
			ReleaseName:    "argocd",
			Chart:          "argo/argo-cd",
			Namespace:      "argocd",
			ValuesTemplate: "templates/argocd-values.yaml.tmpl",
		},
	)

	return &Scenario{
		Steps: steps,
	}
}
