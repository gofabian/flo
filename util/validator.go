package util

import (
	"fmt"
)

type Validator struct {
	Errors []error
}

func (v *Validator) Validate(condition bool, msg string, msgArgs ...interface{}) {
	if !condition {
		v.Error(fmt.Errorf(msg, msgArgs...))
	}
}

func (v *Validator) Error(err error) {
	v.Errors = append(v.Errors, err)
}

func (v *Validator) ValidateEquals(value interface{}, expected interface{}, name string) {
	v.Validate(value == expected, "Expected %s='%s', but found %s='%s'", name, expected, name, value)
}
