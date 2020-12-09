package concourse

import (
	"fmt"
	"os"
	"text/template"
)

type repository struct {
	Branches []branch
}

type branch struct {
	Name           string
	HarmonizedName string
}

func CreateRepositoryPipeline(templateName string, branches []string, file *os.File) error {
	cfg := repository{}
	cfg.Branches = make([]branch, len(branches))

	for i, branch := range branches {
		cfg.Branches[i].Name = branch
		cfg.Branches[i].HarmonizedName = harmonizeName(branch)
	}

	t, err := template.New(templateName).Parse(repositoryPipelineTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse template %s: %w", templateName, err)
	}

	err = t.Execute(file, cfg)
	if err != nil {
		return fmt.Errorf("cannot execute template %s: %w", templateName, err)
	}
	return nil
}
