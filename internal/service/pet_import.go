package service

import (
	"fmt"

	"petsmanagement/internal/model"
)

func (s *Service) ImportPetDraft(key, name string, validationEnabled bool) (err error) {
	draft := model.PetDraft{Key: key, Name: name}
	s.store.SavePetDraft(draft)
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("pet draft validation failed: %v", recovered)
		}
	}()
	validator := model.OptionalPetDraftValidator(validationEnabled)
	if validator != nil {
		if err := validator.Validate(draft); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) PetDraft(key string) (model.PetDraft, bool) { return s.store.PetDraft(key) }
