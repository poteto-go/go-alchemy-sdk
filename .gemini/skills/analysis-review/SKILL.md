---
name: analysis-review
description: Use after a code review to build an action plan for fixing issues. Required when score ≤ 95.
---

# Analysis Review

## Overview

After receiving a code review result, build an **action plan** to fix the issues. Always address the most critical issues first.

**Pass threshold: > 95 points.** If the score is 95 or below, an action plan is **mandatory**.

## Rules

1. **Prioritize by severity** — Fix the highest-deduction issues first.
2. **Ignore typos** — Typos do not require fixes.
3. **Stop when > 95** — Once the projected score exceeds 95, the plan is complete.

## Procedure

1. Read the review result.
2. Sort deductions by severity (largest first).
3. Discard typos.
4. Build a numbered action plan with expected score recovery.
5. Verify projected final score > 95.

---

## Example

### Review Result (Input)

```
## Code Review: wallet/transfer.go

**Score: 62/100**

### Deductions

| # | Category | Detail | Deduction |
|---|----------|--------|-----------|
| 1 | Architecture | Bidirectional dep: wallet ↔ ledger | -20 |
| 2 | Architecture | Duplicated validation in Transfer and Withdraw | -10 |
| 3 | Unit Test | WhiteBox test: TestTransfer_internalLedger | -3 |
| 4 | Unit Test | Asserts mock call args: TestTransfer | -3 |
| 5 | Best Practice | Error not wrapped in Transfer | -3 |
| 6 | Unit Test | Uncovered lines 45-49 (5 lines) | -5 |
| 7 | Best Practice | Typo in comment: "recieve" | 0 |

**Total Deductions: -38** (typo excluded from scoring but noted)
```

### Action Plan (Output)

```
## Action Plan: wallet/transfer.go

**Current Score: 62/100**
**Target: > 95**

### Plan (by priority)

| # | Action | Recovers | Projected |
|---|--------|----------|-----------|
| 1 | Break bidirectional dep: extract shared types to `wallet/types` package, remove `wallet → ledger → wallet` cycle | +20 | 82 |
| 2 | Extract duplicated validation into `validateTransfer()` shared by `Transfer` and `Withdraw` | +10 | 92 |
| 3 | Rewrite TestTransfer_internalLedger as BlackBox test in `wallet_test` package using public API | +3 | 95 (still ≤ 95) |
| 4 | Rewrite TestTransfer to assert balances instead of mock call arguments | +3 | 98 ✅ |

**Projected Score: 98/100** ✅

### Not Fixed
- Typo "recieve" → ignored (typos are not fixed)
- Error wrapping (-3) → not required to reach > 95, fix if time permits
- Uncovered lines 45-49 (-5) → not required to reach > 95, fix if time permits
```

---

## Action Plan Format

```
## Action Plan: [component/file]

**Current Score: XX/100**
**Target: > 95**

### Plan (by priority)

| # | Action | Recovers | Projected |
|---|--------|----------|-----------|
| 1 | [highest severity fix] | +N | XX |
| 2 | [next severity fix] | +N | XX |
| ... | ... | ... | ... |

**Projected Score: XX/100** ✅ or ❌

### Not Fixed
- [items not addressed and why]
```

### Rules for the Plan

- **Order**: Largest deduction first, always.
- **Stop condition**: Once projected score > 95, remaining items go to "Not Fixed".
- **Typos**: Always go to "Not Fixed" — never plan a typo fix.
- **Projected column**: Running total after each fix is applied.
- **If > 95 is unreachable**: Plan all fixes, mark ❌, and note the gap.
