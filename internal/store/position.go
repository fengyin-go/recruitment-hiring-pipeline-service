package store

import (
	"time"

	"recruit/internal/model"
)

func (s *MemoryStore) CreatePosition(p *model.Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.positions {
		if exist.Title == p.Title && exist.DepartmentID == p.DepartmentID {
			return ErrConflict
		}
	}
	s.positions[p.ID] = p
	return nil
}

func (s *MemoryStore) GetPosition(id string) (*model.Position, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.positions[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (s *MemoryStore) ListPositions() []*model.Position {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Position, 0, len(s.positions))
	for _, p := range s.positions {
		cp := *p
		list = append(list, &cp)
	}
	return list
}

// IncrementPositionHiredCount 原子地累加职位的已入职人数，避免并发入职时读改写丢失更新。
func (s *MemoryStore) IncrementPositionHiredCount(id string, delta int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.positions[id]
	if !ok {
		return ErrNotFound
	}
	p.HiredCount += delta
	p.UpdatedAt = time.Now()
	return nil
}

func (s *MemoryStore) UpdatePosition(p *model.Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.positions[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.positions {
		if exist.ID != p.ID && exist.Title == p.Title && exist.DepartmentID == p.DepartmentID {
			return ErrConflict
		}
	}
	s.positions[p.ID] = p
	return nil
}

func (s *MemoryStore) DeletePosition(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.positions[id]; !ok {
		return ErrNotFound
	}
	delete(s.positions, id)
	return nil
}
