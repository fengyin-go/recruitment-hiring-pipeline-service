package service

import (
	"fmt"
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/pkg/idgen"
)

func (s *Service) CreateDepartment(input model.Department) (*model.Department, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d := &model.Department{ID: idgen.Hex(), Name: input.Name, CreatedAt: time.Now()}
	if err := s.store.CreateDepartment(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDepartment(id string) (*model.Department, error) {
	v, err := s.store.GetDepartment(id)
	if err != nil {
		return nil, fmt.Errorf("get department: %v", err)
	}
	return v, nil
}

func (s *Service) ListDepartments(page, size int) ([]*model.Department, int, error) {
	all := s.store.ListDepartments()
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Department{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (s *Service) UpdateDepartment(id string, input model.Department) (*model.Department, error) {
	existing, err := s.store.GetDepartment(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		existing.Name = input.Name
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDepartment(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteDepartment(id string) error {
	// 跨实体校验：部门下仍有职位时不允许删除
	for _, p := range s.store.ListPositions() {
		if p.DepartmentID == id {
			return model.NewValidationError("department", "该部门下仍有职位，无法删除")
		}
	}
	return s.store.DeleteDepartment(id)
}
