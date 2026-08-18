# BUG_REPRO

## Bug 是什么

并发接受多个 Offer 让候选人入职时，职位（Position）上的「已入职人数 HiredCount」会丢失更新：最终统计出来的已入职人数比实际入职的候选人少一两个。与此同时，并发读取职位统计时存在 data race。

## 如何触发

1. 创建一个职位（headcount 足够大，例如 100）。
2. 创建 20 个候选人，并依次走完「安排面试 → 面试通过 → 发放 Offer」流程，使他们都处于 `offered` 状态。
3. 用 20 个 goroutine 在同一时刻同时调用 `AcceptOffer` 接受这些 Offer（可同时并发读取职位统计）。
4. 最终读取该职位的 `HiredCount`，会发现它小于 20。

用下面的命令可以稳定复现（20 次里每次都触发 data race / 计数错误）：

```bash
go test -race -run TestConcurrentHire -count=20 ./internal/service
```

## 错误信息

`-race` 会报告 `Position.HiredCount` 上的 data race，并且测试断言失败：

```
WARNING: DATA RACE
Write at 0x... by goroutine ...:
  recruit/internal/service.(*Service).transitionCandidate()
      internal/service/candidate.go:83

Previous read/write at 0x... by goroutine ...:
  ...

--- FAIL: TestConcurrentHire
    concurrency_test.go:66: hired count = 19, want 20
FAIL
```
