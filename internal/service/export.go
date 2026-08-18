package service

// PositionExportItem 职位导出条目。
type PositionExportItem struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	DepartmentName string `json:"department_name"`
	Headcount      int    `json:"headcount"`
	HiredCount     int    `json:"hired_count"`
	Status         string `json:"status"`
}

// ExportPositions 导出职位汇总。
func (s *Service) ExportPositions() ([]PositionExportItem, error) {
	positions := s.store.ListPositions()
	result := make([]PositionExportItem, 0, len(positions))
	for _, p := range positions {
		item := PositionExportItem{
			ID:         p.ID,
			Title:      p.Title,
			Headcount:  p.Headcount,
			HiredCount: p.HiredCount,
			Status:     p.Status,
		}
		if d, err := s.store.GetDepartment(p.DepartmentID); err == nil {
			item.DepartmentName = d.Name
		}
		result = append(result, item)
	}
	return result, nil
}

// CandidateExportItem 候选人导出条目。
type CandidateExportItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	PositionTitle string `json:"position_title"`
	Status        string `json:"status"`
}

// ExportCandidates 导出候选人汇总。
func (s *Service) ExportCandidates() ([]CandidateExportItem, error) {
	candidates := s.store.ListCandidates()
	result := make([]CandidateExportItem, 0, len(candidates))
	for _, c := range candidates {
		item := CandidateExportItem{
			ID:     c.ID,
			Name:   c.Name,
			Email:  c.Email,
			Status: c.Status,
		}
		if p, err := s.store.GetPosition(c.PositionID); err == nil {
			item.PositionTitle = p.Title
		}
		result = append(result, item)
	}
	return result, nil
}
