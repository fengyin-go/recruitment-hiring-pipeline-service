package model

import (
	"time"
)

const (
	OfferPending  = "pending"
	OfferAccepted = "accepted"
	OfferDeclined = "declined"
)

// offerTransitions Offer 状态机。
var offerTransitions = map[string]map[string]bool{
	OfferPending: {OfferAccepted: true, OfferDeclined: true},
}

// CanTransitionOffer 判断 Offer 状态流转是否合法。
func CanTransitionOffer(from, to string) bool {
	if m, ok := offerTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Offer 录用通知。Salary 单位为「元/月」。
type Offer struct {
	ID          string     `json:"id"`
	CandidateID string     `json:"candidate_id"`
	PositionID  string     `json:"position_id"`
	Salary      int64      `json:"salary"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	DecidedAt   *time.Time `json:"decided_at,omitempty"`
}

func (o *Offer) Validate() error {
	if o.CandidateID == "" {
		return NewValidationError("candidate_id", "候选人不能为空")
	}
	if o.PositionID == "" {
		return NewValidationError("position_id", "职位不能为空")
	}
	if o.Salary <= 0 {
		return NewValidationError("salary", "薪资必须大于 0")
	}
	if o.Status == "" {
		o.Status = OfferPending
	}
	if o.Status != OfferPending && o.Status != OfferAccepted && o.Status != OfferDeclined {
		return NewValidationError("status", "Offer 状态不合法")
	}
	return nil
}

type OfferFilter struct {
	CandidateID string
	PositionID  string
	Status      string
}

func (f OfferFilter) Match(o *Offer) bool {
	if f.CandidateID != "" && o.CandidateID != f.CandidateID {
		return false
	}
	if f.PositionID != "" && o.PositionID != f.PositionID {
		return false
	}
	if f.Status != "" && o.Status != f.Status {
		return false
	}
	return true
}
