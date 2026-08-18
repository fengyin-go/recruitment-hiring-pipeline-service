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
	cp := *p
	s.positions[p.ID] = &cp
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
	cp := *p
	s.positions[p.ID] = &cp
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

// IncrementPositionHiredCount 在写锁内原子地完成「读取职位→HiredCount+1→写回」，
// 避免并发入职时多个 goroutine 各自读出旧值再写回导致的丢失更新。
func (s *MemoryStore) IncrementPositionHiredCount(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.positions[id]
	if !ok {
		return ErrNotFound
	}
	p.HiredCount++
	p.UpdatedAt = time.Now()
	return nil
}
