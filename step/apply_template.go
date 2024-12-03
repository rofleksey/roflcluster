package step

import (
	"github.com/melbahja/goph"
	"roflcluster/config"
)

type ApplyTemplateStep struct {
	File string `yaml:"file"`
}

func (s *ApplyTemplateStep) Execute(client *goph.Client, cfg *config.Config) error {
	data := TemplateData{
		Node: cfg.MainNode,
		User: cfg.User,
	}

	return FormatApplyTemplate(client, s.File, &data)
}

func (s *ApplyTemplateStep) String() string {
	return "Apply template " + s.File
}
