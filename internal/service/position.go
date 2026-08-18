package service

import (
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreatePosition(input model.Position) (*model.Position, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDepartment(input.DepartmentID); err != nil {
		return nil, err
	}
	now := time.Now()
	p := &model.Position{
		ID:           idgen.Hex(),
		Title:        input.Title,
		DepartmentID: input.DepartmentID,
		Description:  input.Description,
		Headcount:    input.Headcount,
		SalaryMin:    input.SalaryMin,
		SalaryMax:    input.SalaryMax,
		Status:       model.PositionOpen,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.CreatePosition(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetPosition(id string) (*model.Position, error) {
	return s.store.GetPosition(id)
}

func (s *Service) ListPositions(filter model.PositionFilter, page, size int) ([]*model.Position, int, error) {
	all := s.store.ListPositions()
	matched := make([]*model.Position, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Position{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdatePosition(id string, input model.Position) (*model.Position, error) {
	existing, err := s.store.GetPosition(id)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		existing.Title = input.Title
	}
	if input.Description != "" {
		existing.Description = input.Description
	}
	if input.Headcount > 0 {
		existing.Headcount = input.Headcount
	}
	if input.SalaryMin >= 0 {
		existing.SalaryMin = input.SalaryMin
	}
	if input.SalaryMax >= 0 {
		existing.SalaryMax = input.SalaryMax
	}
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdatePosition(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeletePosition(id string) error {
	// 跨实体校验：职位被候选人/面试引用时不允许删除
	for _, c := range s.store.ListCandidates() {
		if c.PositionID == id && !model.IsTerminalCandidate(c.Status) {
			return model.NewValidationError("position", "该职位下仍有进行中的候选人，无法删除")
		}
	}
	return s.store.DeletePosition(id)
}
