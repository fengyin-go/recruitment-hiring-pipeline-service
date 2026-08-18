package model

import (
	"strings"
	"time"
)

// Department 招聘部门。
type Department struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *Department) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	if d.Name == "" {
		return NewValidationError("name", "部门名称不能为空")
	}
	return nil
}
