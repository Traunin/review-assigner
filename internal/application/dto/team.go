package dto

type CreateTeamCmd struct {
	TeamName string
	Members  []TeamMemberCmd
}

type TeamMemberCmd struct {
	UserID   string
	Username string
	IsActive bool
}

type TeamDTO struct {
	ID       int64
	TeamName string
}
