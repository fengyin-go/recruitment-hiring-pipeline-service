package service

import (
	"fmt"
	"sort"
	"time"

	"recruit/internal/model"
	"recruit/internal/store"
	"recruit/pkg/idgen"
)

// ScheduleInterview 安排面试：候选人进入面试环节。
func (s *Service) ScheduleInterview(candidateID, positionID, interviewer string, scheduledAt time.Time) (*model.Interview, error) {
	c, err := s.store.GetCandidate(candidateID)
	if err != nil {
		return nil, err
	}
	if c.PositionID != positionID {
		return nil, model.NewValidationError("position_id", "职位与候选人不匹配")
	}
	if model.IsTerminalCandidate(c.Status) {
		return nil, model.NewValidationError("status", "候选人已终态，无法安排面试")
	}
	// 首次进入面试环节时流转状态
	if c.Status == model.CandidateApplied {
		if _, err := s.transitionCandidate(candidateID, model.CandidateInterviewing); err != nil {
			return nil, err
		}
	}
	iv := &model.Interview{
		ID:          idgen.Hex(),
		CandidateID: candidateID,
		PositionID:  positionID,
		Interviewer: interviewer,
		ScheduledAt: scheduledAt,
		Status:      model.InterviewScheduled,
		CreatedAt:   time.Now(),
	}
	if err := iv.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateInterview(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// CompleteInterview 完成面试：通过则进入 Offer 环节，否则淘汰。
func (s *Service) CompleteInterview(id string, passed bool, feedback string) (*model.Interview, error) {
	iv, err := s.store.GetInterview(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionInterview(iv.Status, model.InterviewCompleted) {
		return nil, store.ErrConflict
	}
	iv.Status = model.InterviewCompleted
	iv.Passed = passed
	iv.Feedback = feedback
	if err := s.store.UpdateInterview(iv); err != nil {
		return nil, err
	}
	// 同步候选人状态
	if passed {
		if _, err := s.transitionCandidate(iv.CandidateID, model.CandidateOffered); err != nil {
			return nil, err
		}
	} else {
		if _, err := s.transitionCandidate(iv.CandidateID, model.CandidateRejected); err != nil {
			return nil, err
		}
	}
	return iv, nil
}

func (s *Service) GetInterview(id string) (*model.Interview, error) {
	v, err := s.store.GetInterview(id)
	if err != nil {
		return nil, fmt.Errorf("get interview: %v", err)
	}
	return v, nil
}

func (s *Service) ListInterviews(filter model.InterviewFilter, page, size int) ([]*model.Interview, int, error) {
	all := s.store.ListInterviews()
	matched := make([]*model.Interview, 0, len(all))
	for _, iv := range all {
		if filter.Match(iv) {
			matched = append(matched, iv)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].ScheduledAt.After(matched[j].ScheduledAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Interview{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
