package entities

import "slices"

type Team struct {
	id      TeamID
	name    string
	members []UserID
}

func NewTeam(name string) (*Team, error) {
	if name == "" {
		return nil, ErrTeamNoName
	}
	return &Team{
		name:    name,
		members: make([]UserID, 0),
	}, nil
}

func (t *Team) AddMember(id UserID) error {
	if slices.Contains(t.members, id) {
		return ErrTeamPresent
	}

	t.members = append(t.members, id)
	return nil
}

func (t *Team) RemoveMember(id UserID) {
	t.members = slices.DeleteFunc(t.members, func(member UserID) bool {
		return member == id
	})
}

func (t *Team) ID() TeamID {
	return t.id
}

func (t *Team) Name() string {
	return t.name
}

func (t *Team) Members() []UserID {
	return slices.Clone(t.members)
}
