package store

import (
	"recruit/internal/model"
)

func (s *MemoryStore) CreateReferral(r *model.Referral) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.referrals {
		if exist.CandidateID == r.CandidateID {
			return ErrConflict
		}
	}
	s.referrals[r.ID] = r
	return nil
}

func (s *MemoryStore) GetReferral(id string) (*model.Referral, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.referrals[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListReferrals() []*model.Referral {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Referral, 0, len(s.referrals))
	for _, r := range s.referrals {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateReferral(r *model.Referral) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.referrals[r.ID]; !ok {
		return ErrNotFound
	}
	s.referrals[r.ID] = r
	return nil
}
