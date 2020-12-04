package drone

import (
	"fmt"
	"regexp"

	"github.com/gofabian/flo/validator"
)

var nameRegexp, _ = regexp.Compile("^[a-zA-Z0-9_-]+$")

func ValidatePipeline(p *Pipeline) []error {
	v := &validator.Validator{}

	// global pipeline fields
	v.Validate(p.Kind == "pipeline", "Expected .kind = 'pipeline', but found '%s'", p.Kind)
	v.Validate(p.Type == "docker", "Expected .kind = 'docker', but found '%s'", p.Type)
	v.Validate(nameRegexp.MatchString(p.Name), "Expected .name matches '%s', but found '%s'",
		nameRegexp.String(), p.Name)
	v.Validate(len(p.Steps) > 0, "Expected at least one step at .steps, but found empty array")

	// pipeline steps
	for i, s := range p.Steps {
		v.Validate(nameRegexp.MatchString(s.Name),
			"Expected .steps[%d].name matches '%s', but found '%s'", i, nameRegexp, s.Name)
		v.Validate(len(s.Image) > 0,
			"Expected .steps[%d].image to be non-empty", i)
		for j, cmd := range s.Commands {
			v.Validate(len(cmd) > 0, "Expected .steps[%d].commands[%d] to be non-empty", i, j)
		}

		// step settings
		s.Settings.Decode()
		for _, key := range s.Settings.InvalidKeys {
			v.Error(fmt.Errorf("Expected .steps[%d].settings[%s] to contain string or []string", i, key))
		}
		for key, values := range s.Settings.Fields {
			v.Validate(len(values) > 0, "Expected .steps[%d].settings[%s] to be non-empty", i, key)
		}
	}
	return v.Errors
}
