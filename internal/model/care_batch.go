package model

// CareBatch is a saved group of care labels collected for one pet.
type CareBatch struct {
	PetID  string   `json:"pet_id"`
	Labels []string `json:"labels"`
}

// NewCareBatch prepares a batch for persistence.
func NewCareBatch(petID string, labels []string) *CareBatch {
	return &CareBatch{PetID: petID, Labels: labels}
}
