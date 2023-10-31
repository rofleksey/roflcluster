package step

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/melbahja/goph"
	"gopkg.in/yaml.v2"
	"os"
	"sneakerdocs/config"
	"text/template"
)

type ScenarioStep interface {
	Execute(client *goph.Client, cfg *config.Config) error
	String() string
}

type ScenarioTemp struct {
	Steps []yaml.MapSlice `yaml:"steps"`
}

type Header struct {
	Type string `yaml:"type"`
}

type Scenario struct {
	Steps []ScenarioStep `yaml:"steps"`
}

func unmarshalScenarioStep(stepType string, stepBytes []byte) (ScenarioStep, error) {
	if stepType == "installK3s" {
		return &InstallK3SStep{}, nil
	}

	if stepType == "k3sConfig" {
		return &DownloadK3SConfigStep{}, nil
	}

	if stepType == "tlsSecret" {
		return &TLSSecretStep{}, nil
	}

	if stepType == "applyTemplate" {
		var result ApplyTemplateStep

		err := yaml.Unmarshal(stepBytes, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal step: %w", err)
		}

		return &result, nil
	}

	if stepType == "applyHelm" {
		var result ApplyHelmStep

		err := yaml.Unmarshal(stepBytes, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal step: %w", err)
		}

		return &result, nil
	}

	if stepType == "healthCheck" {
		var result HealthCheckStep

		err := yaml.Unmarshal(stepBytes, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal step: %w", err)
		}

		return &result, nil
	}

	if stepType == "runKubectl" {
		var result RunKubectlStep

		err := yaml.Unmarshal(stepBytes, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal step: %w", err)
		}

		return &result, nil
	}

	return nil, fmt.Errorf("invalid step type: %s", stepType)
}

func ReadScenario(cfg *config.Config) (*Scenario, error) {
	scenarioBytes, err := os.ReadFile("scenario.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read scenario file: %w", err)
	}

	tmpl, err := template.New("scenario").Funcs(sprig.FuncMap()).Parse(string(scenarioBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create scenario template: %w", err)
	}

	var scenarioBuffer bytes.Buffer

	data := TemplateData{
		Node: cfg.MainNode,
		User: cfg.User,
	}

	err = tmpl.Execute(&scenarioBuffer, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to render scenario template: %w", err)
	}

	scenarioBytes = scenarioBuffer.Bytes()

	tempResult := ScenarioTemp{}

	err = yaml.Unmarshal(scenarioBytes, &tempResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal scenario file: %w", err)
	}

	var result Scenario

	for _, tempStep := range tempResult.Steps {
		var header Header

		stepBytes, err := yaml.Marshal(tempStep)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal step type: %w", err)
		}

		err = yaml.Unmarshal(stepBytes, &header)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal step type: %w", err)
		}

		step, err := unmarshalScenarioStep(header.Type, stepBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal step: %w", err)
		}

		result.Steps = append(result.Steps, step)
	}

	return &result, nil
}
