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
	s.offers[o.ID] = o
	return nil
}

func (s *MemoryStore) GetOffer(id string) (*model.Offer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.offers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o, nil
}

func (s *MemoryStore) ListOffers() []*model.Offer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Offer, 0, len(s.offers))
	for _, o := range s.offers {
		list = append(list, o)
	}
	return list
}

func (s *MemoryStore) UpdateOffer(o *model.Offer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.offers[o.ID]; !ok {
		return ErrNotFound
	}
	s.offers[o.ID] = o
	return nil
}
