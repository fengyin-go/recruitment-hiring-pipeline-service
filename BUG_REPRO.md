# BUG_REPRO

## Bug 是什么

通过 HTTP 接口查询一个不存在的实体（部门 / 职位 / 候选人 / 简历 / 面试 / Offer / 内推）时，返回的是 500（服务器内部错误），而预期应该是 404（不存在）。原因是服务层在透传数据访问层的 `not found` 错误时，用 `%v` 格式动词重新包了一层，破坏了错误链，导致处理层用 `errors.Is` 判断「记录不存在」时永远为 false，于是落到了「服务器内部错误」分支。

## 如何触发

启动服务后（或直接跑测试），对任意一个查询接口传入一个不存在的 ID，例如：

```bash
curl -s -i http://localhost:8080/api/departments/not-exist
```

返回状态码是 `500`，响应体类似 `{"code":500,"message":"get department: 记录不存在"}`。

用测试复现（会失败）：

```bash
go test -run TestGetMissingEntityReturnsNotFound ./internal/handler
```

## 错误信息

```
--- FAIL: TestGetMissingEntityReturnsNotFound
    notfound_test.go:22: /api/departments/missing status = 500, want 404
FAIL
```
