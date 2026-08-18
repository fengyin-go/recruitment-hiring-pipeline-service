package store

import (
	"recruit/internal/model"
)

func (s *MemoryStore) CreateOffer(o *model.Offer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.offers {
		if exist.CandidateID == o.CandidateID {
			return ErrConflict
		}
	}
	cp := *o
	s.offers[o.ID] = &cp
	return nil
}

func (s *MemoryStore) GetOffer(id string) (*model.Offer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.offers[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *o
	return &cp, nil
}

func (s *MemoryStore) ListOffers() []*model.Offer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Offer, 0, len(s.offers))
	for _, o := range s.offers {
		cp := *o
		list = append(list, &cp)
	}
	return list
}

func (s *MemoryStore) UpdateOffer(o *model.Offer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.offers[o.ID]; !ok {
		return ErrNotFound
	}
	cp := *o
	s.offers[o.ID] = &cp
	return nil
}
