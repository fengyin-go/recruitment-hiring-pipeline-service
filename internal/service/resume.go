package service

import (
	"fmt"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreateResume(input model.Resume) (*model.Resume, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCandidate(input.CandidateID); err != nil {
		return nil, err
	}
	now := time.Now()
	r := &model.Resume{
		ID:          idgen.Hex(),
		CandidateID: input.CandidateID,
		Summary:     input.Summary,
		Skills:      input.Skills,
		Experience:  input.Experience,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateResume(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetResume(id string) (*model.Resume, error) {
	v, err := s.store.GetResume(id)
	if err != nil {
		return nil, fmt.Errorf("get resume: %v", err)
	}
	return v, nil
}

func (s *Service) GetResumeByCandidate(candidateID string) (*model.Resume, error) {
	v, err := s.store.GetResumeByCandidate(candidateID)
	if err != nil {
		return nil, fmt.Errorf("get resume by candidate: %v", err)
	}
	return v, nil
}

func (s *Service) UpdateResume(id string, input model.Resume) (*model.Resume, error) {
	existing, err := s.store.GetResume(id)
	if err != nil {
		return nil, err
	}
	if input.Summary != "" {
		existing.Summary = input.Summary
	}
	if input.Skills != nil {
		existing.Skills = input.Skills
	}
	if input.Experience != "" {
		existing.Experience = input.Experience
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateResume(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
