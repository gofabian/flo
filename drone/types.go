package drone

import "gopkg.in/yaml.v3"

type Pipeline struct {
	Kind  string
	Type  string
	Name  string
	Steps []Step
}

type Step struct {
	Name     string
	Image    string
	Commands []string
	Settings StepSettings
}

type StepSettings struct {
	RawFields   map[string]yaml.Node `yaml:"-,inline"`
	Fields      map[string][]string  `yaml:"-"`
	InvalidKeys []string             `yaml:"-"`
}

func (s *StepSettings) Decode() {
	if s.Fields != nil && s.InvalidKeys != nil {
		return
	}

	s.Fields = make(map[string][]string)
	s.InvalidKeys = make([]string, 0, len(s.RawFields))
	for key, value := range s.RawFields {
		multi, err := s.decodeField(value)
		if err == nil {
			s.Fields[key] = multi
		} else {
			i := len(s.InvalidKeys)
			s.InvalidKeys = s.InvalidKeys[0 : i+1]
			s.InvalidKeys[i] = key
		}
	}
}

func (s *StepSettings) decodeField(field yaml.Node) ([]string, error) {
	var multi []string
	err := field.Decode(&multi)
	if err == nil {
		return multi, nil
	}
	var single string
	err = field.Decode(&single)
	if err == nil {
		multi = make([]string, 1)
		multi[0] = single
		return multi, nil
	}
	return nil, err
}
