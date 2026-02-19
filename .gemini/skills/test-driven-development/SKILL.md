---
name: test-driven-development
description: Use when implementing any feature or bugfix, before writing implementation code
---

# Test-Driven Development (TDD)

## Overview

Write the test first. Watch it fail. Write minimal code to pass.

**Core principle:** If you didn't watch the test fail, you don't know if it tests the right thing.

**Violating the letter of the rules is violating the spirit of the rules.**

## When to Use

**Always:**

- New features
- Bug fixes
- Refactoring
- Behavior changes

**Exceptions (ask your human partner):**

- Throwaway prototypes
- Generated code
- Configuration files

Thinking "skip TDD just this once"? Stop. That's rationalization.

## The Iron Law

```
NO PRODUCTION CODE WITHOUT A FAILING TEST FIRST
```

Write code before the test? Delete it. Start over.

**No exceptions:**

- Don't keep it as "reference"
- Don't "adapt" it while writing tests
- Don't look at it
- Delete means delete

Implement fresh from tests. Period.

## Good Tests

| Quality          | Good                                | Bad                                                 |
| ---------------- | ----------------------------------- | --------------------------------------------------- |
| **BlackBox**     | Name describes behavior and easy for domain experts to understand             | describes implementation                                     |
| **Minimal**      | One thing. "and" in name? Split it. | two things in one test |
| **Minimal Mocking** | No mocks unless unavoidable except that external services are mocked            | Mocking internal dependencies                        |
| **Refactoring Tolerance**| Assert result, not implementation| Assert details of implementation |


## Red-Green-Refactor

1. RED: Write failing test
2. Watch it fail
3. GREEN: Write minimal code to pass
4. Watch it pass
5. GREEN: Refactor
6. Repeat

### RED - Write Failing Test

Write one minimal test showing what should happen.

```go title="calculator_test.go"
...
func TestSumAbsOfTwoNumbers(t *testing.T) {
    // Act
    result := calculator.AddAbs(2, 3)
    
    // Assert
    assert.Equal(t, 5, result)
}
```

**Requirements:**

- One behavior
- Real code (no mocks unless unavoidable)

### Verify RED - Watch It Fail

**MANDATORY. Never skip.**

```bash
just ut ./calculator/...
```

Confirm:

- Test fails (not errors)

**Test passes?** You're testing existing behavior. Fix test.


### GREEN - Minimal Code

```go title="calculator.go"
...
func AddAbs(a, b int) int {
    return a + b
}
```

### Verify GREEN - Watch It Pass

**MANDATORY.**

```bash
just ut ./calculator/...
```

Confirm:

- Test passes
- Other tests still pass
- Output pristine (no errors, warnings)

**Test fails?** Fix code, not test.


### GREEN - Refactor

After green only:

- Remove duplication
- Improve names
- Extract helpers

Keep tests green. Don't add behavior.

```go title="calculator_test.go"
...
func TestSumAbsOfTwoNumbers(t *testing.T) {
    // Act & Assert
    assert.Equal(t, 5, calculator.AddAbs(2, 3))
}
```

### Repeat

Next failing test for next feature.

one more example cycle for `AddAbs`:

**RED Failing Test:**

```go title="calculator_test.go"
...
func TestSumAbsOfTwoNumbers(t *testing.T) {
    t.Run("sum of two positive numbers", func(t *testing.T) {
        // Act & Assert
        assert.Equal(t, 5, calculator.AddAbs(2, 3))
    })

    t.Run("sum of two negative numbers", func(t *testing.T) {
        // Act & Assert
        assert.Equal(t, 5, calculator.AddAbs(-2, -3))
    })
}
```

test failed.

**GREEN Code:**

```go title="calculator.go"
...
func AddAbs(a, b int) int {
    if a < 0 {
        a = -a
    }
    if b < 0 {
        b = -b
    }
    return a + b
}
```

test passed.

**GREEN Refactor:**

```go title="calculator.go"
...
func abs(a int) int {
  if a < 0 {
    return -a
  }
  return a
}

func AddAbs(a, b int) int {
    return abs(a) + abs(b)
}
```

test passed.

**FINISH:**

If there's no more to implement & do refactoring, you're done.


## Why Order Matters

**"I'll write tests after to verify it works"**

Tests written after code pass immediately. Passing immediately proves nothing:

- Might test wrong thing
- Might test implementation, not behavior
- Might miss edge cases you forgot
- You never saw it catch the bug

Test-first forces you to see the test fail, proving it actually tests something.

**"I already manually tested all the edge cases"**

Manual testing is ad-hoc. You think you tested everything but:

- No record of what you tested
- Can't re-run when code changes
- Easy to forget cases under pressure
- "It worked when I tried it" ≠ comprehensive

Automated tests are systematic. They run the same way every time.

**"Deleting X hours of work is wasteful"**

Sunk cost fallacy. The time is already gone. Your choice now:

- Delete and rewrite with TDD (X more hours, high confidence)
- Keep it and add tests after (30 min, low confidence, likely bugs)

The "waste" is keeping code you can't trust. Working code without real tests is technical debt.

**"TDD is dogmatic, being pragmatic means adapting"**

TDD IS pragmatic:

- Finds bugs before commit (faster than debugging after)
- Prevents regressions (tests catch breaks immediately)
- Documents behavior (tests show how to use code)
- Enables refactoring (change freely, tests catch breaks)

"Pragmatic" shortcuts = debugging in production = slower.

**"Tests after achieve the same goals - it's spirit not ritual"**

No. Tests-after answer "What does this do?" Tests-first answer "What should this do?"

Tests-after are biased by your implementation. You test what you built, not what's required. You verify remembered edge cases, not discovered ones.

Tests-first force edge case discovery before implementing. Tests-after verify you remembered everything (you didn't).

30 minutes of tests after ≠ TDD. You get coverage, lose proof tests work.

## Common Rationalizations

| Excuse                                 | Reality                                                                 |
| -------------------------------------- | ----------------------------------------------------------------------- |
| "Too simple to test"                   | Simple code breaks. Test takes 30 seconds.                              |
| "I'll test after"                      | Tests passing immediately prove nothing.                                |
| "Tests after achieve same goals"       | Tests-after = "what does this do?" Tests-first = "what should this do?" |
| "Already manually tested"              | Ad-hoc ≠ systematic. No record, can't re-run.                           |
| "Deleting X hours is wasteful"         | Sunk cost fallacy. Keeping unverified code is technical debt.           |
| "Keep as reference, write tests first" | You'll adapt it. That's testing after. Delete means delete.             |
| "Need to explore first"                | Fine. Throw away exploration, start with TDD.                           |
| "Test hard = design unclear"           | Listen to test. Hard to test = hard to use.                             |
| "TDD will slow me down"                | TDD faster than debugging. Pragmatic = test-first.                      |
| "Manual test faster"                   | Manual doesn't prove edge cases. You'll re-test every change.           |
| "Existing code has no tests"           | You're improving it. Add tests for existing code.                       |

## Red Flags - STOP and Start Over

- Code before test
- Test after implementation
- Test passes immediately
- Can't explain why test failed
- Tests added "later"
- Rationalizing "just this once"
- "I already manually tested it"
- "Tests after achieve the same purpose"
- "It's about spirit not ritual"
- "Keep as reference" or "adapt existing code"
- "Already spent X hours, deleting is wasteful"
- "TDD is dogmatic, I'm being pragmatic"
- "This is different because..."

**All of these mean: Delete code. Start over with TDD.**

## Example: Bug Fix

**Bug:** Empty email accepted

**RED**

```go
func Test_SendMail(t *testing.T) {
    t.Run("rejects empty email", func(t *testing.T) {
        _, err := SendMail("")
        assert.Error(t, err)
    })
}
```

**Verify RED**

```bash
$ just ut
FAIL: expected error, got nil
```

**GREEN**

```go
func SendMail(email string) (bool, error) {
    if email == "" {
        return false, errors.New("empty email")
    }

    ...
}
```

**Verify GREEN**

```bash
$ just ut
PASS
```

**REFACTOR**
Extract validation for multiple fields if needed.

## Verification Checklist

Before marking work complete:

- [ ] Every new function/method has a test
- [ ] Watched each test fail before implementing
- [ ] Each test failed for expected reason (feature missing, not typo)
- [ ] Wrote minimal code to pass each test
- [ ] All tests pass
- [ ] Output pristine (no errors, warnings)
- [ ] Tests use real code (mocks only if unavoidable)
- [ ] Edge cases and errors covered

Can't check all boxes? You skipped TDD. Start over.

## When Stuck

| Problem                | Solution                                                             |
| ---------------------- | -------------------------------------------------------------------- |
| Don't know how to test | Write wished-for API. Write assertion first. Ask your human partner. |
| Test too complicated   | Design too complicated. Simplify interface.                          |
| Must mock everything   | Code too coupled. Use dependency injection.                          |
| Test setup huge        | Extract helpers. Still complex? Simplify design.                     |

