package service

import (
	"fmt"
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

// CreateReferral 登记内推：校验候选人与职位存在。
func (s *Service) CreateReferral(referrer, candidateID, positionID string, bonus int64) (*model.Referral, error) {
	c, err := s.store.GetCandidate(candidateID)
	if err != nil {
		return nil, err
	}
	if c.PositionID != positionID {
		return nil, model.NewValidationError("position_id", "职位与候选人不匹配")
	}
	r := &model.Referral{
		ID:          idgen.Hex(),
		Referrer:    referrer,
		CandidateID: candidateID,
		PositionID:  positionID,
		Bonus:       bonus,
		Status:      model.ReferralPending,
		CreatedAt:   time.Now(),
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateReferral(r); err != nil {
		return nil, err
	}
	return r, nil
}

// MarkReferralHired 内推成功（候选人入职后调用）。
func (s *Service) MarkReferralHired(id string) (*model.Referral, error) {
	r, err := s.store.GetReferral(id)
	if err != nil {
		return nil, err
	}
	if r.Status == model.ReferralHired {
		return r, nil
	}
	now := time.Now()
	r.Status = model.ReferralHired
	r.HiredAt = &now
	if err := s.store.UpdateReferral(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetReferral(id string) (*model.Referral, error) {
	v, err := s.store.GetReferral(id)
	if err != nil {
		return nil, fmt.Errorf("get referral: %v", err)
	}
	return v, nil
}

func (s *Service) ListReferrals(filter model.ReferralFilter, page, size int) ([]*model.Referral, int, error) {
	all := s.store.ListReferrals()
	matched := make([]*model.Referral, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Referral{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
