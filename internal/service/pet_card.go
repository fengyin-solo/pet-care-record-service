package service

import (
	"fmt"

	"petsmanagement/internal/model"
)

func (s *Service) BuildPetCard(key, name string, sections []string) (card *model.PetCard, err error) {
	card = &model.PetCard{Key: key, Name: name}
	s.store.PutPetCard(card)
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("build pet card: %v", recovered)
		}
	}()
	for _, section := range sections {
		card.AddSection(section)
	}
	card.Ready = true
	return card, nil
}

func (s *Service) PetCard(key string) (*model.PetCard, bool) { return s.store.PetCard(key) }
