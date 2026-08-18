// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"recruit/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	// 部门
	CreateDepartment(d *model.Department) error
	GetDepartment(id string) (*model.Department, error)
	ListDepartments() []*model.Department
	UpdateDepartment(d *model.Department) error
	DeleteDepartment(id string) error

	// 职位
	CreatePosition(p *model.Position) error
	GetPosition(id string) (*model.Position, error)
	ListPositions() []*model.Position
	UpdatePosition(p *model.Position) error
	DeletePosition(id string) error
	IncrementPositionHiredCount(id string, delta int) error

	// 候选人
	CreateCandidate(c *model.Candidate) error
	GetCandidate(id string) (*model.Candidate, error)
	ListCandidates() []*model.Candidate
	UpdateCandidate(c *model.Candidate) error
	DeleteCandidate(id string) error

	// 简历
	CreateResume(r *model.Resume) error
	GetResume(id string) (*model.Resume, error)
	GetResumeByCandidate(candidateID string) (*model.Resume, error)
	ListResumes() []*model.Resume
	UpdateResume(r *model.Resume) error

	// 面试
	CreateInterview(i *model.Interview) error
	GetInterview(id string) (*model.Interview, error)
	ListInterviews() []*model.Interview
	UpdateInterview(i *model.Interview) error

	// Offer
	CreateOffer(o *model.Offer) error
	GetOffer(id string) (*model.Offer, error)
	ListOffers() []*model.Offer
	UpdateOffer(o *model.Offer) error

	// 内推
	CreateReferral(r *model.Referral) error
	GetReferral(id string) (*model.Referral, error)
	ListReferrals() []*model.Referral
	UpdateReferral(r *model.Referral) error
}
