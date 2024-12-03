package step

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/melbahja/goph"
	"os"
	"path/filepath"
	"roflcluster/config"
	"strings"
	"text/template"
)

type TemplateData struct {
	Node config.NodeConfig
	User config.UserConfig
}

func FormatUploadTemplate(client *goph.Client, templatePath string, data any) error {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(sprig.FuncMap()).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	_ = os.MkdirAll("temp", os.ModePerm)

	renderPath := filepath.Join("temp", strings.TrimSuffix(filepath.Base(templatePath), ".tmpl"))
	defer func() {
		os.Remove(renderPath)
	}()

	file, err := os.Create(renderPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", renderPath, err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to render template %s: %w", templatePath, err)
	}

	err = client.Upload(renderPath, filepath.Base(renderPath))
	if err != nil {
		return fmt.Errorf("failed to upload rendered template %s: %w", renderPath, err)
	}

	return nil
}

func FormatApplyTemplate(client *goph.Client, templatePath string, data any) error {
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(sprig.FuncMap()).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	_ = os.MkdirAll("temp", os.ModePerm)

	renderPath := filepath.Join("temp", strings.TrimSuffix(filepath.Base(templatePath), ".tmpl"))
	defer func() {
		os.Remove(renderPath)
	}()

	file, err := os.Create(renderPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", renderPath, err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to render template %s: %w", templatePath, err)
	}

	err = client.Upload(renderPath, filepath.Base(renderPath))
	if err != nil {
		return fmt.Errorf("failed to upload rendered template %s: %w", renderPath, err)
	}

	installCmd := fmt.Sprintf(fmt.Sprintf("kubectl --kubeconfig /etc/rancher/k3s/k3s.yaml apply -f %s", filepath.Base(renderPath)))
	stdout, err := client.Run(installCmd)
	if err != nil {
		return fmt.Errorf("failed to run apply template %s: %w %s", renderPath, err, string(stdout))
	}

	return nil
}
