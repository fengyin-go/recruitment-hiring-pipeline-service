package store

import (
	"recruit/internal/model"
)

func (s *MemoryStore) CreateCandidate(c *model.Candidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.candidates {
		if exist.Email == c.Email && exist.PositionID == c.PositionID {
			return ErrConflict
		}
	}
	s.candidates[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCandidate(id string) (*model.Candidate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.candidates[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListCandidates() []*model.Candidate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Candidate, 0, len(s.candidates))
	for _, c := range s.candidates {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCandidate(c *model.Candidate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.candidates[c.ID]; !ok {
		return ErrNotFound
	}
	s.candidates[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCandidate(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.candidates[id]; !ok {
		return ErrNotFound
	}
	delete(s.candidates, id)
	return nil
}
