package model

import (
	"strings"
	"time"
)

const (
	CandidateApplied      = "applied"
	CandidateInterviewing = "interviewing"
	CandidateOffered      = "offered"
	CandidateHired        = "hired"
	CandidateRejected     = "rejected"
)

// candidateTransitions 候选人状态机。
var candidateTransitions = map[string]map[string]bool{
	CandidateApplied:      {CandidateInterviewing: true, CandidateRejected: true},
	CandidateInterviewing: {CandidateOffered: true, CandidateRejected: true},
	CandidateOffered:      {CandidateHired: true, CandidateRejected: true},
}

// CanTransitionCandidate 判断候选人状态流转是否合法。
func CanTransitionCandidate(from, to string) bool {
	if m, ok := candidateTransitions[from]; ok {
		return m[to]
	}
	return false
}

// IsTerminalCandidate 判断是否为终态。
func IsTerminalCandidate(status string) bool {
	return status == CandidateHired || status == CandidateRejected
}

// Candidate 候选人。
type Candidate struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	PositionID string     `json:"position_id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

func (c *Candidate) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	if c.Name == "" {
		return NewValidationError("name", "候选人姓名不能为空")
	}
	if c.Email == "" {
		return NewValidationError("email", "候选人邮箱不能为空")
	}
	if c.PositionID == "" {
		return NewValidationError("position_id", "应聘职位不能为空")
	}
	if c.Status == "" {
		c.Status = CandidateApplied
	}
	if c.Status != CandidateApplied && c.Status != CandidateInterviewing &&
		c.Status != CandidateOffered && c.Status != CandidateHired && c.Status != CandidateRejected {
		return NewValidationError("status", "候选人状态不合法")
	}
	return nil
}

type CandidateFilter struct {
	PositionID string
	Status     string
	Keyword    string
}

func (f CandidateFilter) Match(c *Candidate) bool {
	if f.PositionID != "" && c.PositionID != f.PositionID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) &&
			!strings.Contains(strings.ToLower(c.Email), k) {
			return false
		}
	}
	return true
}
