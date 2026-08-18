package service

import (
	"sort"

	"recruit/internal/model"
)

// KeyCount 通用键值计数。
type KeyCount struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// CandidateStats 候选人统计。
type CandidateStats struct {
	Total      int        `json:"total"`
	ByStatus   []KeyCount `json:"by_status"`
	ByPosition []KeyCount `json:"by_position"`
}

// CandidateStats 计算候选人统计。
func (s *Service) CandidateStats() (*CandidateStats, error) {
	candidates := s.store.ListCandidates()
	stats := &CandidateStats{Total: len(candidates)}
	statusMap := map[string]int{}
	posMap := map[string]int{}
	for _, c := range candidates {
		statusMap[c.Status]++
		posMap[c.PositionID]++
	}
	stats.ByStatus = keyCounts(statusMap, func(k string) string { return k })
	stats.ByPosition = keyCounts(posMap, func(k string) string {
		if p, err := s.store.GetPosition(k); err == nil {
			return p.Title
		}
		return k
	})
	return stats, nil
}

// FunnelStage 招聘漏斗阶段。
type FunnelStage struct {
	Stage string `json:"stage"`
	Count int    `json:"count"`
}

// RecruitmentFunnel 招聘漏斗统计（投递→面试→Offer→入职）。
func (s *Service) RecruitmentFunnel() (*[]FunnelStage, error) {
	candidates := s.store.ListCandidates()
	applied := 0
	interviewing := 0
	offered := 0
	hired := 0
	for _, c := range candidates {
		switch c.Status {
		case model.CandidateApplied:
			applied++
		case model.CandidateInterviewing:
			interviewing++
		case model.CandidateOffered:
			offered++
		case model.CandidateHired:
			hired++
		}
	}
	// 漏斗各阶段为累积值（面试过的包含后续阶段）
	reachedInterview := interviewing + offered + hired
	reachedOffer := offered + hired
	funnel := []FunnelStage{
		{Stage: "applied", Count: applied + reachedInterview},
		{Stage: "interview", Count: reachedInterview},
		{Stage: "offer", Count: reachedOffer},
		{Stage: "hired", Count: hired},
	}
	return &funnel, nil
}

// PositionStats 职位招聘进度统计。
type PositionStats struct {
	PositionID    string `json:"position_id"`
	PositionTitle string `json:"position_title"`
	DepartmentID  string `json:"department_id"`
	Headcount     int    `json:"headcount"`
	HiredCount    int    `json:"hired_count"`
	Candidates    int    `json:"candidates"`
}

// ListPositionStats 统计每个职位的招聘进度。
func (s *Service) ListPositionStats() ([]PositionStats, error) {
	positions := s.store.ListPositions()
	result := make([]PositionStats, 0, len(positions))
	for _, p := range positions {
		ps := PositionStats{
			PositionID:    p.ID,
			PositionTitle: p.Title,
			DepartmentID:  p.DepartmentID,
			Headcount:     p.Headcount,
			HiredCount:    p.HiredCount,
		}
		for _, c := range s.store.ListCandidates() {
			if c.PositionID == p.ID {
				ps.Candidates++
			}
		}
		result = append(result, ps)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Candidates > result[j].Candidates })
	return result, nil
}

// InterviewStats 面试统计（通过率）。
type InterviewStats struct {
	Total     int     `json:"total"`
	Completed int     `json:"completed"`
	Passed    int     `json:"passed"`
	PassRate  float64 `json:"pass_rate"`
}

// InterviewStats 计算面试通过率。
func (s *Service) InterviewStats() (*InterviewStats, error) {
	interviews := s.store.ListInterviews()
	stats := &InterviewStats{Total: len(interviews)}
	for _, iv := range interviews {
		if iv.Status == model.InterviewCompleted {
			stats.Completed++
			if iv.Passed {
				stats.Passed++
			}
		}
	}
	if stats.Completed > 0 {
		stats.PassRate = float64(stats.Passed) / float64(stats.Completed)
	}
	return stats, nil
}

// DepartmentStats 按部门统计职位与候选人数量。
type DepartmentStats struct {
	DepartmentID   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
	Positions      int    `json:"positions"`
	Candidates     int    `json:"candidates"`
}

// ListDepartmentStats 统计各部门招聘数据。
func (s *Service) ListDepartmentStats() ([]DepartmentStats, error) {
	departments := s.store.ListDepartments()
	result := make([]DepartmentStats, 0, len(departments))
	for _, d := range departments {
		ds := DepartmentStats{DepartmentID: d.ID, DepartmentName: d.Name}
		posSet := map[string]bool{}
		for _, p := range s.store.ListPositions() {
			if p.DepartmentID == d.ID {
				ds.Positions++
				posSet[p.ID] = true
			}
		}
		for _, c := range s.store.ListCandidates() {
			if posSet[c.PositionID] {
				ds.Candidates++
			}
		}
		result = append(result, ds)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Candidates > result[j].Candidates })
	return result, nil
}

// keyCounts 将 map 计数转为排序后的 KeyCount 列表。
func keyCounts(m map[string]int, name func(string) string) []KeyCount {
	result := make([]KeyCount, 0, len(m))
	for k, c := range m {
		result = append(result, KeyCount{Key: name(k), Count: c})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Count > result[j].Count })
	return result
}

// MonthlyCount 按月聚合的投递数量。
type MonthlyCount struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

// MonthlyStats 按月统计候选人投递数量。
func (s *Service) MonthlyStats() ([]MonthlyCount, error) {
	candidates := s.store.ListCandidates()
	countMap := map[string]int{}
	for _, c := range candidates {
		key := c.CreatedAt.Format("2006-01")
		countMap[key]++
	}
	result := make([]MonthlyCount, 0, len(countMap))
	for m, c := range countMap {
		result = append(result, MonthlyCount{Month: m, Count: c})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Month > result[j].Month })
	return result, nil
}

// ReferralStats 内推统计。
type ReferralStats struct {
	Total      int   `json:"total"`
	Hired      int   `json:"hired"`
	TotalBonus int64 `json:"total_bonus"`
}

// ReferralStats 计算内推数量与奖励总额。
func (s *Service) ReferralStats() (*ReferralStats, error) {
	referrals := s.store.ListReferrals()
	stats := &ReferralStats{Total: len(referrals)}
	for _, r := range referrals {
		if r.Status == model.ReferralHired {
			stats.Hired++
			stats.TotalBonus += r.Bonus
		}
	}
	return stats, nil
}
