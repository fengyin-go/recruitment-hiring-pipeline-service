package service

import (
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/internal/store"
	"recruit/pkg/idgen"
)

// CreateOffer 发放 Offer：候选人进入待接受状态。
func (s *Service) CreateOffer(candidateID, positionID string, salary int64) (*model.Offer, error) {
	c, err := s.store.GetCandidate(candidateID)
	if err != nil {
		return nil, err
	}
	if c.PositionID != positionID {
		return nil, model.NewValidationError("position_id", "职位与候选人不匹配")
	}
	if c.Status != model.CandidateOffered {
		if !model.CanTransitionCandidate(c.Status, model.CandidateOffered) {
			return nil, store.ErrConflict
		}
		if _, err := s.transitionCandidate(candidateID, model.CandidateOffered); err != nil {
			return nil, err
		}
	}
	o := &model.Offer{
		ID:          idgen.Hex(),
		CandidateID: candidateID,
		PositionID:  positionID,
		Salary:      salary,
		Status:      model.OfferPending,
		CreatedAt:   time.Now(),
	}
	if err := o.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateOffer(o); err != nil {
		return nil, err
	}
	return o, nil
}

// AcceptOffer 候选人接受 Offer：入职。
func (s *Service) AcceptOffer(id string) (*model.Offer, error) {
	o, err := s.store.GetOffer(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionOffer(o.Status, model.OfferAccepted) {
		return nil, store.ErrConflict
	}
	now := time.Now()
	o.Status = model.OfferAccepted
	o.DecidedAt = &now
	if err := s.store.UpdateOffer(o); err != nil {
		return nil, err
	}
	if _, err := s.transitionCandidate(o.CandidateID, model.CandidateHired); err != nil {
		return nil, err
	}
	return o, nil
}

// DeclineOffer 候选人拒绝 Offer：淘汰。
func (s *Service) DeclineOffer(id string) (*model.Offer, error) {
	o, err := s.store.GetOffer(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionOffer(o.Status, model.OfferDeclined) {
		return nil, store.ErrConflict
	}
	now := time.Now()
	o.Status = model.OfferDeclined
	o.DecidedAt = &now
	if err := s.store.UpdateOffer(o); err != nil {
		return nil, err
	}
	if _, err := s.transitionCandidate(o.CandidateID, model.CandidateRejected); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) GetOffer(id string) (*model.Offer, error) {
	return s.store.GetOffer(id)
}

func (s *Service) ListOffers(filter model.OfferFilter, page, size int) ([]*model.Offer, int, error) {
	all := s.store.ListOffers()
	matched := make([]*model.Offer, 0, len(all))
	for _, o := range all {
		if filter.Match(o) {
			matched = append(matched, o)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Offer{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
