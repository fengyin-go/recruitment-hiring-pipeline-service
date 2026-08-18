package store

import (
	"recruit/internal/model"
)

// cloneResume 深拷贝 Resume：Skills 为切片，浅拷贝会与调用方共享底层数组。
func cloneResume(r *model.Resume) *model.Resume {
	cp := *r
	if r.Skills != nil {
		cp.Skills = append([]string(nil), r.Skills...)
	}
	return &cp
}

func (s *MemoryStore) CreateResume(r *model.Resume) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.resumes {
		if exist.CandidateID == r.CandidateID {
			return ErrConflict
		}
	}
	s.resumes[r.ID] = cloneResume(r)
	return nil
}

func (s *MemoryStore) GetResume(id string) (*model.Resume, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.resumes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneResume(r), nil
}

func (s *MemoryStore) GetResumeByCandidate(candidateID string) (*model.Resume, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.resumes {
		if r.CandidateID == candidateID {
			return cloneResume(r), nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListResumes() []*model.Resume {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Resume, 0, len(s.resumes))
	for _, r := range s.resumes {
		list = append(list, cloneResume(r))
	}
	return list
}

func (s *MemoryStore) UpdateResume(r *model.Resume) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.resumes[r.ID]; !ok {
		return ErrNotFound
	}
	s.resumes[r.ID] = cloneResume(r)
	return nil
}
