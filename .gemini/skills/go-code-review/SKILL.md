---
name: go-code-review
description: Use when reviewing Go code changes. Scores 0-100 with deductions per category.
---

# Go Code Review

## Overview

Score every code review from **100 → 0**. Start at 100 and apply deductions.

Issues detectable by linters (staticcheck, golangci-lint, etc.) are **NOT** scored here.

## Scoring Summary

| Category | Deduction | Per |
|---|---|---|
| Bidirectional dependency | -20 | per occurrence |
| Duplicated logic | -10 | per occurrence |
| Violates project design (AGENTS.md) | -10 | per occurrence |
| Non-refactoring-tolerant test | -3 | per test case |
| WhiteBox test | -3 | per test case |
| Test package not `xxx_test` | -3 | per file |
| Uncovered code | -1 | per line |
| Security HIGH | -10 | per finding |
| Security MEDIUM | -5 | per finding |
| Security LOW | -1 | per finding |
| Spec not met | -100 | — |
| Go best practice violation | varies | per occurrence |

---

## Architecture

### Bidirectional Dependency (-20 per occurrence)

Packages must depend in one direction only. Types-only packages (e.g. `types`, `model`) are allowed as shared dependencies.

```go
// ❌ BAD: bidirectional dependency
// package order imports package payment
package order

import "myapp/payment"

func NewOrder(p payment.Processor) { ... }

// package payment imports package order
package payment

import "myapp/order"

func Charge(o order.Order) { ... }
```

```go
// ✅ GOOD: unidirectional with shared types
// package types (shared, no imports from order or payment)
package types

type Order struct {
    ID     string
    Amount int
}

// package payment depends on types only
package payment

import "myapp/types"

func Charge(o types.Order) error { ... }

// package order depends on types and payment
package order

import (
    "myapp/types"
    "myapp/payment"
)

func Place(o types.Order) error {
    return payment.Charge(o)
}
```

### Duplicated Logic (-10 per occurrence)

Extract shared behavior into a function or package.

```go
// ❌ BAD: duplicated validation in two handlers
func CreateUser(name string) error {
    if len(name) == 0 || len(name) > 100 {
        return errors.New("invalid name")
    }
    // ...
}

func UpdateUser(name string) error {
    if len(name) == 0 || len(name) > 100 {
        return errors.New("invalid name")
    }
    // ...
}
```

```go
// ✅ GOOD: shared validation
func validateName(name string) error {
    if len(name) == 0 || len(name) > 100 {
        return errors.New("invalid name")
    }
    return nil
}

func CreateUser(name string) error {
    if err := validateName(name); err != nil {
        return err
    }
    // ...
}

func UpdateUser(name string) error {
    if err := validateName(name); err != nil {
        return err
    }
    // ...
}
```

### Project Design Violation (-10 per occurrence)

Code must follow the architecture described in the project's `AGENTS.md`. Check layer boundaries, naming conventions, and package responsibilities.

---

## Unit Testing

### Non-Refactoring-Tolerant Test (-3 per test case)

Tests must assert **results**, not **how** the code works internally. Asserting arguments, call order, or internal state couples tests to implementation.

```go
// ❌ BAD: asserts arguments passed to internal dependency
func TestTransfer(t *testing.T) {
    mockRepo := new(MockRepo)
    svc := NewService(mockRepo)

    svc.Transfer("alice", "bob", 100)

    // Asserting internal call arguments = refactoring breaks this
    mockRepo.AssertCalledWith(t, "Debit", "alice", 100)
    mockRepo.AssertCalledWith(t, "Credit", "bob", 100)
}
```

```go
// ✅ GOOD: asserts observable outcome
func TestTransfer(t *testing.T) {
    repo := NewInMemoryRepo()
    repo.SetBalance("alice", 200)
    repo.SetBalance("bob", 50)
    svc := NewService(repo)

    err := svc.Transfer("alice", "bob", 100)

    assert.NoError(t, err)
    assert.Equal(t, 100, repo.GetBalance("alice"))
    assert.Equal(t, 150, repo.GetBalance("bob"))
}
```

### WhiteBox Test (-3 per test case)

Always use **BlackBox testing**. Tests must exercise the public API from an external perspective.

```go
// ❌ BAD: WhiteBox — test in same package, accesses unexported internals
package wallet

func TestWallet_internalLedger(t *testing.T) {
    w := &wallet{ledger: []entry{}}
    w.addEntry(100) // unexported method
    assert.Equal(t, 1, len(w.ledger))
}
```

```go
// ✅ GOOD: BlackBox — test in _test package, uses only exported API
package wallet_test

import (
    "testing"
    "myapp/wallet"
    "github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
    w := wallet.New()

    err := w.Deposit(100)

    assert.NoError(t, err)
    assert.Equal(t, 100, w.Balance())
}
```

### Test Package Naming (-3 per file)

Test files **must** use the `xxx_test` package, not `xxx`.

```go
// ❌ BAD
package wallet

// ✅ GOOD
package wallet_test
```

### Uncovered Code (-1 per line)

Every line of production code should be exercised by tests. Run:

```bash
just coverage ./...
```

---

## Security

Review security against the project context (README.md / AGENTS.md).

| Severity | Deduction | Example |
|---|---|---|
| **HIGH** (-10) | SQL injection, hardcoded secrets, missing auth | `db.Exec("SELECT * FROM users WHERE id=" + id)` |
| **MEDIUM** (-5) | Insufficient input validation, weak crypto | Using `md5` for password hashing |
| **LOW** (-1) | Missing rate limiting, verbose error messages | Returning internal error details to client |

```go
// ❌ HIGH: SQL injection
func GetUser(db *sql.DB, id string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = '" + id + "'")
    // ...
}
```

```go
// ✅ GOOD: parameterized query
func GetUser(db *sql.DB, id string) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
    // ...
}
```

```go
// ❌ MEDIUM: secrets in code
const apiKey = "sk-live-abc123xyz"
```

```go
// ✅ GOOD: secrets from environment
apiKey := os.Getenv("API_KEY")
```

---

## Specification

### Spec Not Met (-100)

If the code does not satisfy the specification, the review score is effectively 0. Verify behavior against the requirements before anything else.

---

## Go Best Practices

### Error Handling (-5 per occurrence)

```go
// ❌ BAD: silently ignoring error
data, _ := json.Marshal(obj)
```

```go
// ✅ GOOD: handle or propagate
data, err := json.Marshal(obj)
if err != nil {
    return fmt.Errorf("marshal obj: %w", err)
}
```

### Error Wrapping (-3 per occurrence)

```go
// ❌ BAD: context lost
if err != nil {
    return err
}
```

```go
// ✅ GOOD: wrap with context
if err != nil {
    return fmt.Errorf("create user %q: %w", name, err)
}
```

### Interface Design (-5 per occurrence)

Define interfaces where they are **used**, not where they are implemented. Keep interfaces small.

```go
// ❌ BAD: large interface defined in implementation package
package storage

type Storage interface {
    Get(key string) ([]byte, error)
    Set(key string, val []byte) error
    Delete(key string) error
    List(prefix string) ([]string, error)
    Watch(key string) <-chan Event
    Backup() error
    Restore(path string) error
}
```

```go
// ✅ GOOD: small interface defined at consumer
package order

type ItemStore interface {
    Get(key string) ([]byte, error)
    Set(key string, val []byte) error
}

type Service struct {
    store ItemStore
}
```

### Context Propagation (-3 per occurrence)

```go
// ❌ BAD: creating new context inside
func FetchData() ([]byte, error) {
    ctx := context.Background()
    return client.Get(ctx, "/data")
}
```

```go
// ✅ GOOD: accept context from caller
func FetchData(ctx context.Context) ([]byte, error) {
    return client.Get(ctx, "/data")
}
```

### Goroutine Leak (-5 per occurrence)

```go
// ❌ BAD: goroutine never stops
func StartWorker() {
    go func() {
        for {
            process()
        }
    }()
}
```

```go
// ✅ GOOD: cancellable via context
func StartWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                process()
            }
        }
    }()
}
```

### Naked Returns (-2 per occurrence)

```go
// ❌ BAD: hard to read, especially in long functions
func parse(s string) (result int, err error) {
    // ... 30 lines ...
    return
}
```

```go
// ✅ GOOD: explicit return values
func parse(s string) (int, error) {
    // ...
    return result, nil
}
```

### Package Naming (-2 per occurrence)

```go
// ❌ BAD
package util     // too generic
package helpers  // too generic
package common   // too generic
```

```go
// ✅ GOOD
package auth     // describes purpose
package invoice  // describes domain
```

### Struct Initialization (-2 per occurrence)

```go
// ❌ BAD: positional fields — brittle if struct changes
u := User{"alice", 30, true}
```

```go
// ✅ GOOD: named fields
u := User{
    Name:   "alice",
    Age:    30,
    Active: true,
}
```

---

## Review Procedure

1. **Spec** — Does the code meet the requirements? If not, stop (-100).
2. **Architecture** — Check dependency direction, duplication, project design.
3. **Unit Tests** — BlackBox? Refactoring-tolerant? `_test` package? Coverage?
4. **Security** — Review against project context for vulnerabilities.
5. **Go Best Practices** — Error handling, interfaces, context, goroutines, naming.
6. **Calculate Score** — Start at 100, apply all deductions. Minimum is 0.

### Output Format

```
## Code Review: [component/file]

**Score: XX/100**

### Deductions

| # | Category | Detail | Deduction |
|---|----------|--------|-----------|
| 1 | Architecture | Bidirectional dep: order ↔ payment | -20 |
| 2 | Unit Test | WhiteBox test: TestWallet_internal | -3 |
| ... | ... | ... | ... |

**Total Deductions: -XX**

### Summary
[brief summary of key issues and recommendations]
```
