package mapper

import (
	"github.com/Traunin/review-assigner/internal/application/dto"
	"github.com/Traunin/review-assigner/internal/domain/entities"
)

func ToReviewerDTO(r entities.Reviewer) dto.ReviewerDTO {
	return dto.ReviewerDTO{
		UserID:     r.UserID,
		AssignedAt: r.AssignedAt,
	}
}

func ToReviewerDTOs(domain []entities.Reviewer) []dto.ReviewerDTO {
	out := make([]dto.ReviewerDTO, len(domain))
	for i, r := range domain {
		out[i] = ToReviewerDTO(r)
	}
	return out
}

func ToPullRequestDTO(pr *entities.PullRequest) dto.PullRequestDTO {
	mergedAt := pr.MergedAt()

	return dto.PullRequestDTO{
		PullRequestID:   pr.ID(),
		PullRequestName: pr.Name(),
		AuthorID:        pr.AuthorID(),
		Status:          string(pr.Status()),
		CreatedAt:       pr.CreatedAt(),
		MergedAt:        &mergedAt,
		Reviewers:       ToReviewerDTOs(pr.Reviewers()),
	}
}

func ToUserDTO(u *entities.User) dto.UserDTO {
	var tid int64
	if u.TeamID() != nil {
		tid = int64(*u.TeamID())
	}
	return dto.UserDTO{
		UserID:   u.ID(),
		Username: u.Username(),
		IsActive: u.IsActive(),
		TeamID:   &tid,
	}
}

func ToTeamDTO(t *entities.Team, members []entities.User) dto.TeamDTO {
	users := make([]dto.UserDTO, len(members))
	for i := range members {
		users[i] = ToUserDTO(&members[i])
	}

	return dto.TeamDTO{
		ID:       int64(t.ID()),
		TeamName: t.Name(),
	}
}

func ToUserDTOs(domain []*entities.User) []*dto.UserDTO {
	out := make([]*dto.UserDTO, len(domain))
	for i, u := range domain {
		dto := ToUserDTO(u)
		out[i] = &dto
	}
	return out
}
