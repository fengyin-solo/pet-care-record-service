package model

import "errors"

type PetDraft struct {
	Key       string
	Name      string
	Validated bool
}

func (d *PetDraft) MarkValidated() {}

type PetDraftValidator interface {
	Validate(PetDraft) error
}

type requiredNameValidator struct{ prefix *string }

func (v *requiredNameValidator) Validate(draft PetDraft) error {
	if draft.Name == "" || *v.prefix == "" {
		return errors.New("pet name is required")
	}
	return nil
}

func OptionalPetDraftValidator(enabled bool) PetDraftValidator {
	if !enabled {
		var validator *requiredNameValidator
		return validator
	}
	prefix := "pet"
	return &requiredNameValidator{prefix: &prefix}
}
