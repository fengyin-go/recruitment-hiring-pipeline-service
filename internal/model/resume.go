package model

import (
	"strings"
	"time"
)

// Resume 候选人简历。
type Resume struct {
	ID          string    `json:"id"`
	CandidateID string    `json:"candidate_id"`
	Summary     string    `json:"summary"`
	Skills      []string  `json:"skills"`
	Experience  string    `json:"experience"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Resume) Validate() error {
	r.Summary = strings.TrimSpace(r.Summary)
	if r.CandidateID == "" {
		return NewValidationError("candidate_id", "候选人不能为空")
	}
	if r.Summary == "" {
		return NewValidationError("summary", "简历摘要不能为空")
	}
	return nil
}
