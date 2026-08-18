package store

import (
	"recruit/internal/model"
)

func (s *MemoryStore) CreateInterview(i *model.Interview) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interviews[i.ID] = i
	return nil
}

func (s *MemoryStore) GetInterview(id string) (*model.Interview, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.interviews[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) ListInterviews() []*model.Interview {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Interview, 0, len(s.interviews))
	for _, i := range s.interviews {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) UpdateInterview(i *model.Interview) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.interviews[i.ID]; !ok {
		return ErrNotFound
	}
	s.interviews[i.ID] = i
	return nil
}
