package entities

type UserID string
type TeamID int
type PrID string

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

func (s PRStatus) String() string {
	return string(s)
}
