package validator

import (
	"fmt"

	"kcl-lang.io/kcl-go"
)

type Validator struct {
	values  []byte
	schemas []byte
}

func New(values, schemas []byte) *Validator {
	return &Validator{
		values:  values,
		schemas: schemas,
	}
}

func (v *Validator) Validate() error {
	opts := &kcl.ValidateOptions{Format: "yaml", Schema: "Values"}
	_, err := kcl.ValidateCode(string(v.values), string(v.schemas), opts)
	if err != nil {
		return fmt.Errorf("error validating: %v", err)
	}
	return nil
}
