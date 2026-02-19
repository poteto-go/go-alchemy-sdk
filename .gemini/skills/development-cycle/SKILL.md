---
name: development-cycle
description: Use to follow the standard development workflow. Two modes — normal implementation and optimization.
---

# Development Cycle

## Overview

Two development cycles. Both end with a CI-passing commit.

**If anything in the user's instructions is unclear, ask the user before proceeding.**

---

## Normal Implementation

```
User Instruction → TDD → CI → Review → Plan → (repeat) → CI → Commit
```

### Steps

1. **Implement with TDD** (@test-driven-development)
   - Write failing test → watch it fail → write minimal code → watch it pass → refactor.

2. **Run CI**
   ```bash
   just ci
   ```
   - If CI fails, go back to step 1 and fix.
   - Repeat 1 → 2 until CI passes.

3. **Code Review** (@go-code-review / @rust-code-review / @ts-code-review / @moonbit-code-review)
   - Run the review for the relevant language.
   - Score > 95 → skip to step 6.

4. **Action Plan** (@analysis-review)
   - Score < 95 → build a prioritized action plan.

5. **Fix and Re-review**
   - Go back to step 1 with the action plan.
   - Repeat 1 → 4 until score > 95.

6. **Final CI and Commit**
   ```bash
   just ci
   ```
   - CI passes → commit following @git-commit.

### Flow

```text
       ┌─────────────────────┐
       │                     │
       v                     │
  1. TDD Implement           │
       │                     │
       v                     │
  2. just ci ── FAIL ────────┘
       │
      PASS
       │
       v
  3. Code Review
       │
       ├── > 95 ────────┐
       │                │
       v                │
  4. Action Plan        │
       │                │
       v                │
  5. Fix (go to 1) ─┐   │
       │             │   │
       v             │   │
    Re-review ≤ 95 ─┘   │
       │                │
      > 95              │
       v                │
  6. just ci + commit <─┘
```

---

## Optimization Implementation

```
User Instruction → Measure → Find Bottleneck → TDD → CI → Review → Plan → (repeat) → CI → Commit
```

### Steps

1. **Find Optimization Opportunity** (@performance-tuning)
   - Measure baseline (build size / execution speed).
   - Profile and identify the biggest bottleneck.
   - Write benchmark tests if needed.

2. **Implement with TDD** (@test-driven-development)
   - Optimize the bottleneck following the TDD cycle.

3. **Run CI**
   ```bash
   just ci
   ```
   - If CI fails, go back to step 2 and fix.
   - Repeat 2 → 3 until CI passes.

4. **Measure After** (@performance-tuning)
   - Run the same measurements as step 1.
   - If improvement < 1% and readability is worse, revert.

5. **Code Review** (@go-code-review / @rust-code-review / @ts-code-review / @moonbit-code-review)
   - Score > 95 → skip to step 8.

6. **Action Plan** (@analysis-review)
   - Score ≤ 95 → build a prioritized action plan.

7. **Fix and Re-review**
   - Go back to step 2 with the action plan.
   - Repeat 2 → 6 until score > 95.

8. **Final CI and Commit**
   ```bash
   just ci
   ```
   - CI passes → commit following @git-commit.
   - Include performance results in the commit message body.
