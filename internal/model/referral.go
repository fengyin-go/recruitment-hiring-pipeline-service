package model

import (
	"strings"
	"time"
)

const (
	ReferralPending = "pending"
	ReferralHired   = "hired"
)

// Referral 内推记录。Bonus 单位为「元」（整数，避免浮点）。
type Referral struct {
	ID          string     `json:"id"`
	Referrer    string     `json:"referrer"`
	CandidateID string     `json:"candidate_id"`
	PositionID  string     `json:"position_id"`
	Bonus       int64      `json:"bonus"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	HiredAt     *time.Time `json:"hired_at,omitempty"`
}

func (r *Referral) Validate() error {
	r.Referrer = strings.TrimSpace(r.Referrer)
	if r.Referrer == "" {
		return NewValidationError("referrer", "推荐人不能为空")
	}
	if r.CandidateID == "" {
		return NewValidationError("candidate_id", "候选人不能为空")
	}
	if r.PositionID == "" {
		return NewValidationError("position_id", "职位不能为空")
	}
	if r.Bonus < 0 {
		return NewValidationError("bonus", "内推奖励不能为负")
	}
	if r.Status == "" {
		r.Status = ReferralPending
	}
	if r.Status != ReferralPending && r.Status != ReferralHired {
		return NewValidationError("status", "内推状态不合法")
	}
	return nil
}

type ReferralFilter struct {
	Referrer string
	Status   string
}

func (f ReferralFilter) Match(r *Referral) bool {
	if f.Referrer != "" && r.Referrer != f.Referrer {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
