# BUG_REPRO

## Bug 是什么

候选人已经被淘汰（状态为 `rejected`）之后，仍然可以对其面试调用「完成面试」，接口会返回错误，但面试记录却被持久化成了「已完成」（`completed`）。也就是说：一个已淘汰的候选人，却留有一条「已完成且通过」的面试记录，状态不一致。

## 如何触发

1. 创建职位 + 候选人。
2. 给候选人安排一场面试（此时面试为 `scheduled`，候选人进入 `interviewing`）。
3. 把候选人淘汰（`RejectCandidate`，候选人变为 `rejected`）。
4. 调用 `CompleteInterview` 完成该面试（传入 `passed=true`）。
5. 观察：方法返回错误，但面试记录的 `Status` 已变为 `completed`。

用测试复现（会失败）：

```bash
go test -run TestCompleteInterviewRejectedCandidateInconsistent ./internal/service
```

## 错误信息

```
--- FAIL: TestCompleteInterviewRejectedCandidateInconsistent
    interview_consistency_test.go:36: interview status = completed, want scheduled (should not be completed)
FAIL
```
