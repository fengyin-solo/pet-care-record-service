package model

type RequestScope struct {
	OwnerID string
	PetID   string
	Labels  []string
	Owned   bool
}

func (s *RequestScope) Prepare(ownerID, petID string, labels []string) {
	s.OwnerID = ownerID
	s.PetID = petID
	s.Labels = labels
}
