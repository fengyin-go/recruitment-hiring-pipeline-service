package store

import (
	"sync"

	"recruit/internal/model"
)

// MemoryStore 基于内存的 Store 实现，使用读写锁保证并发安全。
type MemoryStore struct {
	mu          sync.RWMutex
	departments map[string]*model.Department
	positions   map[string]*model.Position
	candidates  map[string]*model.Candidate
	resumes     map[string]*model.Resume
	interviews  map[string]*model.Interview
	offers      map[string]*model.Offer
	referrals   map[string]*model.Referral
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		departments: make(map[string]*model.Department),
		positions:   make(map[string]*model.Position),
		candidates:  make(map[string]*model.Candidate),
		resumes:     make(map[string]*model.Resume),
		interviews:  make(map[string]*model.Interview),
		offers:      make(map[string]*model.Offer),
		referrals:   make(map[string]*model.Referral),
	}
}

var _ Store = (*MemoryStore)(nil)
