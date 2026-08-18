package service

import (
	"testing"
	"time"

	"recruit/internal/model"
)

// TestCompleteInterviewRejectedCandidateInconsistent 验证：候选人已被淘汰时，完成其面试应当失败，
// 且不应把面试记录标记为已完成。
func TestCompleteInterviewRejectedCandidateInconsistent(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, err := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	if err != nil {
		t.Fatalf("create candidate: %v", err)
	}
	iv, err := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())
	if err != nil {
		t.Fatalf("schedule interview: %v", err)
	}
	// 面试前候选人被淘汰
	if _, err := s.RejectCandidate(c.ID); err != nil {
		t.Fatalf("reject candidate: %v", err)
	}
	// 完成一个已淘汰候选人的面试，应当返回错误
	if _, err := s.CompleteInterview(iv.ID, true, ""); err == nil {
		t.Fatalf("expect error completing interview for rejected candidate")
	}
	// 面试记录不应被标记为已完成
	gotIv, err := s.GetInterview(iv.ID)
	if err != nil {
		t.Fatalf("get interview: %v", err)
	}
	if gotIv.Status != model.InterviewScheduled {
		t.Fatalf("interview status = %s, want scheduled (should not be completed)", gotIv.Status)
	}
}
