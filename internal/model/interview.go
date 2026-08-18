package model

import (
	"strings"
	"time"
)

const (
	InterviewScheduled = "scheduled"
	InterviewCompleted = "completed"
	InterviewCanceled  = "canceled"
)

// interviewTransitions 面试状态机。
var interviewTransitions = map[string]map[string]bool{
	InterviewScheduled: {InterviewCompleted: true, InterviewCanceled: true},
}

// CanTransitionInterview 判断面试状态流转是否合法。
func CanTransitionInterview(from, to string) bool {
	if m, ok := interviewTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Interview 面试安排。
type Interview struct {
	ID          string    `json:"id"`
	CandidateID string    `json:"candidate_id"`
	PositionID  string    `json:"position_id"`
	Interviewer string    `json:"interviewer"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      string    `json:"status"`
	Feedback    string    `json:"feedback"`
	Passed      bool      `json:"passed"`
	CreatedAt   time.Time `json:"created_at"`
}

func (i *Interview) Validate() error {
	i.Interviewer = strings.TrimSpace(i.Interviewer)
	i.Feedback = strings.TrimSpace(i.Feedback)
	if i.CandidateID == "" {
		return NewValidationError("candidate_id", "候选人不能为空")
	}
	if i.PositionID == "" {
		return NewValidationError("position_id", "职位不能为空")
	}
	if i.Interviewer == "" {
		return NewValidationError("interviewer", "面试官不能为空")
	}
	if i.Status == "" {
		i.Status = InterviewScheduled
	}
	if i.Status != InterviewScheduled && i.Status != InterviewCompleted && i.Status != InterviewCanceled {
		return NewValidationError("status", "面试状态不合法")
	}
	return nil
}

type InterviewFilter struct {
	CandidateID string
	PositionID  string
	Status      string
}

func (f InterviewFilter) Match(i *Interview) bool {
	if f.CandidateID != "" && i.CandidateID != f.CandidateID {
		return false
	}
	if f.PositionID != "" && i.PositionID != f.PositionID {
		return false
	}
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	return true
}
