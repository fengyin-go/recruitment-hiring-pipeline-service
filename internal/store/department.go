package store

import (
	"recruit/internal/model"
)

func (s *MemoryStore) CreateDepartment(d *model.Department) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.departments {
		if exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.departments[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDepartment(id string) (*model.Department, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.departments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDepartments() []*model.Department {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Department, 0, len(s.departments))
	for _, d := range s.departments {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDepartment(d *model.Department) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.departments[d.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.departments {
		if exist.ID != d.ID && exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.departments[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDepartment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.departments[id]; !ok {
		return ErrNotFound
	}
	delete(s.departments, id)
	return nil
}
