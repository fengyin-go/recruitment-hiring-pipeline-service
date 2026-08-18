package service

import (
	"fmt"
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/internal/store"
	"recruit/pkg/idgen"
)

// CreateCandidate 候选人投递：校验职位开放并生成候选人记录。
func (s *Service) CreateCandidate(input model.Candidate) (*model.Candidate, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	pos, err := s.store.GetPosition(input.PositionID)
	if err != nil {
		return nil, err
	}
	if pos.Status != model.PositionOpen {
		return nil, model.NewValidationError("position_id", "该职位已停止招聘")
	}
	now := time.Now()
	c := &model.Candidate{
		ID:         idgen.Hex(),
		Name:       input.Name,
		Email:      input.Email,
		Phone:      input.Phone,
		PositionID: input.PositionID,
		Status:     model.CandidateApplied,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateCandidate(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCandidate(id string) (*model.Candidate, error) {
	v, err := s.store.GetCandidate(id)
	if err != nil {
		return nil, fmt.Errorf("get candidate: %v", err)
	}
	return v, nil
}

func (s *Service) ListCandidates(filter model.CandidateFilter, page, size int) ([]*model.Candidate, int, error) {
	all := s.store.ListCandidates()
	matched := make([]*model.Candidate, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Candidate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// transitionCandidate 流转候选人状态（内部）。
func (s *Service) transitionCandidate(id, target string) (*model.Candidate, error) {
	c, err := s.store.GetCandidate(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionCandidate(c.Status, target) {
		return nil, store.ErrConflict
	}
	now := time.Now()
	c.Status = target
	c.UpdatedAt = now
	if model.IsTerminalCandidate(target) {
		c.FinishedAt = &now
		if target == model.CandidateHired {
			// 入职后职位录用数 +1
			if pos, err := s.store.GetPosition(c.PositionID); err == nil {
				pos.HiredCount++
				pos.UpdatedAt = now
				_ = s.store.UpdatePosition(pos)
			}
		}
	}
	if err := s.store.UpdateCandidate(c); err != nil {
		return nil, err
	}
	return c, nil
}

// RejectCandidate 淘汰候选人。
func (s *Service) RejectCandidate(id string) (*model.Candidate, error) {
	return s.transitionCandidate(id, model.CandidateRejected)
}

// BatchReject 批量淘汰候选人，返回成功数。
func (s *Service) BatchReject(ids []string) (int, error) {
	success := 0
	for _, id := range ids {
		if _, err := s.RejectCandidate(id); err == nil {
			success++
		}
	}
	return success, nil
}
