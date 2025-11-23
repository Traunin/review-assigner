package entities

type UserID string
type TeamID int
type PullRequestID string

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

func (s PRStatus) String() string {
	return string(s)
}

func (id UserID) String() string {
	return string(id)
}

func (id PullRequestID) String() string {
	return string(id)
}
