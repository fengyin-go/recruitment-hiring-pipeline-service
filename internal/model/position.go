package model

import (
	"strings"
	"time"
)

const (
	PositionOpen   = "open"
	PositionClosed = "closed"
)

// Position 招聘职位。
// SalaryMin/SalaryMax 单位为「元/月」（整数，避免浮点）。
type Position struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	DepartmentID string    `json:"department_id"`
	Description  string    `json:"description"`
	Headcount    int       `json:"headcount"`
	HiredCount   int       `json:"hired_count"`
	SalaryMin    int64     `json:"salary_min"`
	SalaryMax    int64     `json:"salary_max"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Position) Validate() error {
	p.Title = strings.TrimSpace(p.Title)
	if p.Title == "" {
		return NewValidationError("title", "职位名称不能为空")
	}
	if p.DepartmentID == "" {
		return NewValidationError("department_id", "所属部门不能为空")
	}
	if p.Headcount <= 0 {
		return NewValidationError("headcount", "招聘人数必须大于 0")
	}
	if p.SalaryMin < 0 || p.SalaryMax < 0 {
		return NewValidationError("salary", "薪资不能为负")
	}
	if p.SalaryMax > 0 && p.SalaryMin > p.SalaryMax {
		return NewValidationError("salary", "薪资下限不能高于上限")
	}
	if p.Status == "" {
		p.Status = PositionOpen
	}
	if p.Status != PositionOpen && p.Status != PositionClosed {
		return NewValidationError("status", "职位状态不合法")
	}
	return nil
}

type PositionFilter struct {
	DepartmentID string
	Status       string
	Keyword      string
}

func (f PositionFilter) Match(p *Position) bool {
	if f.DepartmentID != "" && p.DepartmentID != f.DepartmentID {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(p.Title), k) {
			return false
		}
	}
	return true
}
