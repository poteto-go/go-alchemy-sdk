---
name: performance-tuning
description: Use when optimizing build size or execution speed. Always measure before and after.
---

# Performance Tuning

## Overview

Two optimization targets:

1. **Build Size Minimization** — Reduce the compiled binary / bundle size.
2. **Execution Speed Maximization** — Reduce runtime latency, including unit test execution time.

**Core principle:** Never optimize without measuring. Never ship without proving the improvement.

## Procedure

```
Measure Before → Find Bottleneck → Optimize → Measure After → Evaluate
```

1. **Measure Before** — Record baseline metrics.
2. **Analyze** — Profile code and tests to find the bottleneck.
3. **Optimize** — Fix the biggest bottleneck first. Write benchmark tests if needed.
4. **Measure After** — Record post-optimization metrics.
5. **Evaluate** — If improvement < 1%, decide based on readability.

---

## Measure Before (Mandatory)

Always capture baseline metrics before making any changes.

### Build Size

```bash
# Go
go build -o app ./cmd/app && ls -lh app
# Example output: -rwxr-xr-x 1 user user 12M Feb 19 app

# TypeScript (Vite)
npx vite build
# Example output: dist/index.js   245.12 kB │ gzip: 78.34 kB

# Rust
cargo build --release && ls -lh target/release/app
# Example output: -rwxr-xr-x 1 user user 4.2M Feb 19 app
```

### Execution Speed

```bash
# Go — unit test time
go test ./... -count=1 2>&1 | tail -5
# Example output:
# ok   myapp/wallet    0.032s
# ok   myapp/order     0.118s
# ok   myapp/payment   0.540s  ← slowest

# Go — benchmark
go test -bench=. -benchmem ./wallet/...
# Example output:
# BenchmarkTransfer-8   50000   32145 ns/op   4096 B/op   12 allocs/op

# TypeScript
npx vitest --reporter=verbose 2>&1 | tail -10

# Rust
cargo test -- --report-time 2>&1 | tail -10
```

### Record Baseline

```
## Baseline (Before)

| Metric | Value |
|---|---|
| Binary size | 12M |
| Total test time | 0.690s |
| Slowest test suite | payment (0.540s) |
| Benchmark (Transfer) | 32145 ns/op, 4096 B/op |
```

---

## Find the Bottleneck

### Analyze Build Size

```bash
# Go — find what contributes to binary size
go build -o app ./cmd/app
go tool nm -size app | sort -rnk2 | head -20

# Go — check for debug info
go build -ldflags="-s -w" -o app_stripped ./cmd/app
ls -lh app app_stripped

# TypeScript — analyze bundle
npx vite build --report
# or
npx source-map-explorer dist/index.js
```

Look for:
- Unused dependencies
- Large vendored libraries with smaller alternatives
- Debug symbols in release builds
- Duplicated code across packages

### Analyze Execution Speed

```bash
# Go — CPU profile
go test -cpuprofile=cpu.prof -bench=. ./payment/...
go tool pprof -top cpu.prof

# Go — find slow tests
go test ./... -v -count=1 2>&1 | grep -E "^\s*(ok|FAIL|---)" | sort -t$'\t' -k2 -rn

# Go — memory profile
go test -memprofile=mem.prof -bench=. ./payment/...
go tool pprof -top mem.prof
```

Look for:
- Hot functions consuming disproportionate CPU
- Excessive allocations
- Slow test setup (DB connections, large fixtures)
- Unnecessary I/O in unit tests

**Always start with the biggest bottleneck.**

---

## Optimize

### Write Benchmark Tests (if needed)

Before optimizing execution speed, write a benchmark to prove the improvement.

```go
// payment/payment_test.go
func BenchmarkProcessPayment(b *testing.B) {
    svc := NewService(NewInMemoryRepo())
    for i := 0; i < b.N; i++ {
        svc.Process(Payment{Amount: 100, Currency: "USD"})
    }
}
```

```bash
# Run benchmark before optimization
go test -bench=BenchmarkProcessPayment -benchmem -count=5 ./payment/... | tee before.txt
```

### Optimization Priorities

**Always fix the biggest bottleneck first.** Common optimizations by category:

#### Build Size

| Technique | Example |
|---|---|
| Strip debug symbols | `go build -ldflags="-s -w"` |
| Remove unused deps | `go mod tidy`, review `go.sum` |
| Replace heavy deps | Use `encoding/json` instead of large reflection-based libs |
| Tree-shaking | Ensure bundler eliminates dead code |
| Compress assets | gzip / brotli static files |

#### Execution Speed

| Technique | Example |
|---|---|
| Reduce allocations | Reuse buffers, pre-allocate slices with `make([]T, 0, cap)` |
| Avoid reflection | Use code generation or direct field access |
| Batch I/O | Bulk insert instead of row-by-row |
| Cache hot paths | `sync.Pool`, in-memory cache for repeated lookups |
| Parallelize | `t.Parallel()` for independent tests, `errgroup` for concurrent work |
| Reduce test I/O | In-memory repos instead of DB calls in unit tests |

---

## Measure After (Mandatory)

Run the exact same measurements as the baseline.

```bash
# Same commands as "Measure Before"
go build -o app ./cmd/app && ls -lh app
go test ./... -count=1 2>&1 | tail -5
go test -bench=BenchmarkProcessPayment -benchmem -count=5 ./payment/... | tee after.txt
```

### Compare

```bash
# Go — benchstat comparison
go install golang.org/x/perf/cmd/benchstat@latest
benchstat before.txt after.txt
```

---

## Evaluate

### Record Results

```
## Results

| Metric | Before | After | Change |
|---|---|---|---|
| Binary size | 12M | 8.4M | -30% ✅ |
| Total test time | 0.690s | 0.412s | -40% ✅ |
| Slowest suite | payment 0.540s | payment 0.285s | -47% ✅ |
| Benchmark (Transfer) | 32145 ns/op | 28900 ns/op | -10% ✅ |
| Allocations | 4096 B/op | 2048 B/op | -50% ✅ |
```

### Decision Criteria

| Improvement | Action |
|---|---|
| ≥ 1% | **Accept** — commit the optimization. |
| < 1% and **readable** | **Accept** — marginal gain but no readability cost. |
| < 1% and **less readable** | **Reject** — revert. Readability wins over negligible gains. |

### Summary Template

```
## Performance Tuning: [component]

### Target
[Build size / Execution speed / Both]

### Bottleneck
[What was identified as the primary bottleneck]

### Changes
1. [optimization 1 — description]
2. [optimization 2 — description]

### Results
| Metric | Before | After | Change |
|---|---|---|---|
| ... | ... | ... | ... |

### Decision
[Accepted / Rejected with rationale]
```
