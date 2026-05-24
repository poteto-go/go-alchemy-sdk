---
name: golang-tdd
description: Guide for Test-Driven Development (TDD) in Go. Use this skill whenever implementing a new feature, fixing a bug, adding a method, or modifying existing Go code. Trigger on phrases like "implement", "add feature", "fix bug", "write test", "TDD cycle", "red green refactor", or any request to write new Go code. Always apply TDD — write the failing test first, then make it pass.
---

# Golang TDD in go-alchemy-sdk

This project enforces TDD. Every code change starts with a failing test. The cycle is always: **Red → Green → Refactor**.

## TDD Cycle

Repeat this loop until the entire spec is satisfied:

### 1. RED — Write a failing test

Write the next test case that describes one unit of desired behavior. Run it and confirm it fails.

```bash
go test ./... -run TestYourNewTest -v
```

The test must fail before you write implementation. A test that passes without code isn't a TDD test.

### 2. GREEN — Write minimum code to pass, then refactor

Write only enough code to make the test pass — no more.

```bash
go test ./... -run TestYourNewTest -v
```

Once green, refactor freely: clean up duplication, improve naming, extract helpers in test or source code. The tests stay green throughout.

```bash
just ut  # confirm nothing broke after refactor
```

### 3. Back to RED — Next test case

Add the next test case for the next piece of behavior. Loop back to step 1. Keep looping until the spec is fully implemented.

---

## Test Patterns Used in This Project

### Table-driven tests (preferred for multiple cases)

```go
func TestRouter_Find(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        wantCode int
    }{
        {"root path", "/", 200},
        {"not found", "/missing", 404},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // arrange
            // act
            // assert
            assert.Equal(t, tt.wantCode, got)
        })
    }
}
```

### HTTP handler tests (use `net/http/httptest`)

```go
func TestHandler_Get(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/hello", nil)
    w := httptest.NewRecorder()

    ctx := NewContext(w, req, bindings)
    err := handler(ctx)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Async behavior (use `assert.Eventually`)

```go
assert.Eventually(t, func() bool {
    return triggered
}, 1*time.Second, 10*time.Millisecond)
```

### Test naming convention

```
TestTypeName_MethodName        — top-level
t.Run("descriptive scenario")  — subtests
```

White-box tests (same package, e.g., `package takibi`) for internal access.
Black-box tests (e.g., `package takibi_test`) for public API contracts.

---

## Workflow Checklist

When asked to implement anything, follow this sequence:

1. **Understand the interface** — read the relevant `interfaces/` types before writing anything
2. **Loop until the spec is satisfied:**
   - Write the next test case (RED) — run it, confirm it fails
   - Write minimum code to make it pass (GREEN) — run it, confirm it passes
   - Refactor test or source code as needed — run `just ut` to confirm still green
   - Move to the next test case
3. **Update `/docs`** — if the public API changed, update the docs directory

---

## Common Pitfalls

- **Don't write implementation before the test.** If you catch yourself writing both at once, stop — write the test file first, run it to see red, then write the implementation.
- **Don't test implementation details.** Test behavior through the public/package interface, not unexported helpers.
- **Avoid mocking the framework internals.** Tests here use `httptest` and real instances — this project was bitten by mock/real divergence before.
- **`-gcflags=all=-l` is required** for coverage to work correctly with inlined functions. Always use `just ut` or the full test command, not bare `go test ./...`.
