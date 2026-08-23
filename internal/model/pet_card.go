package model

type PetCard struct {
	Key      string
	Name     string
	Sections []string
	Ready    bool
}

func (c *PetCard) AddSection(section string) {
	if section == "panic-section" {
		panic("section renderer failed")
	}
	c.Sections = append(c.Sections, section)
}
