package service

import (
	"sync"
	"testing"
	"time"

	"recruit/internal/model"
)

// TestConcurrentHire 验证并发入职时职位已入职人数能正确累计，且统计读取不被并发写入污染。
func TestConcurrentHire(t *testing.T) {
	s, _, pos := setupPosition(t)
	const n = 20
	offers := make([]string, 0, n)
	for i := 0; i < n; i++ {
		c, err := s.CreateCandidate(model.Candidate{
			Name:       string(rune('A' + i)),
			Email:      string(rune('a'+i)) + "@x.com",
			PositionID: pos.ID,
		})
		if err != nil {
			t.Fatalf("create candidate: %v", err)
		}
		iv, err := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())
		if err != nil {
			t.Fatalf("schedule interview: %v", err)
		}
		if _, err := s.CompleteInterview(iv.ID, true, ""); err != nil {
			t.Fatalf("complete interview: %v", err)
		}
		o, err := s.CreateOffer(c.ID, pos.ID, 30000)
		if err != nil {
			t.Fatalf("create offer: %v", err)
		}
		offers = append(offers, o.ID)
	}

	var wg sync.WaitGroup
	start := make(chan struct{})
	for _, id := range offers {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			<-start
			_, _ = s.AcceptOffer(id)
		}(id)
	}
	// 同时读统计，暴露列表返回内部引用导致的并发读脏数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 200; i++ {
			_, _ = s.ListPositionStats()
		}
	}()
	close(start)
	wg.Wait()

	got, err := s.GetPosition(pos.ID)
	if err != nil {
		t.Fatalf("get position: %v", err)
	}
	if got.HiredCount != n {
		t.Fatalf("hired count = %d, want %d", got.HiredCount, n)
	}
}
