package concourse

import (
	"fmt"
	"io"
	"text/template"
)

type repository struct {
	Branches []branch
}

type branch struct {
	Name           string
	HarmonizedName string
}

func CreateRepositoryPipeline(selfUpdateJob bool, branches []string, writer io.Writer) error {
	var templateName string
	if selfUpdateJob && len(branches) > 0 {
		templateName = "full-pipeline"
	} else if selfUpdateJob {
		templateName = "self-update-pipeline"
	} else if len(branches) > 0 {
		templateName = "build-pipeline"
	} else {
		return fmt.Errorf("missing template")
	}

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

	err = t.Execute(writer, cfg)
	if err != nil {
		return fmt.Errorf("cannot execute template %s: %w", templateName, err)
	}
	return nil
}
