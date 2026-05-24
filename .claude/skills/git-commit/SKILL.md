---
name: git-commit
description: Guidelines for creating standardized commit messages with conventional prefixes and issue tracking.
---

# Commit Rule

## Overview
All commits must follow a strict format to ensure traceability and standardization.

## Rules

1.  **Prefix**: Use a conventional prefix to categorize the commit.
2.  **Issue Linking**:
    -   Extract the issue number from the branch name (e.g., `user/#123/feature-x` -> `#123`).
    -   Include `refs: #<issue-number>` in the message.
3.  **Co-Author Attribution**:
    -   Always include the co-author trailer for the AI assistant.
    -   `Co-Authored-By: gemini-cli <218195315+gemini-cli@users.noreply.github.com>`

## Prefixes

| Prefix | Usage |
|---|---|
| `feat` | New feature or functionality |
| `fix` | Bug fix |
| `refactor` | Code restructuring without behavior change |
| `perf` | Performance improvement |
| `test` | Adding or updating tests |
| `chore` | Maintenance, dependency updates, config changes |
| `mig` | Database or data migration |
| `docs` | Documentation only |
| `ci` | CI/CD pipeline changes |
| `style` | Code formatting (no logic change) |

## Format Template

```bash
git commit -m "<prefix>: <message>. refs: #<issue-number>" -m "Co-Authored-By: gemini-cli <218195315+gemini-cli@users.noreply.github.com>"
```

## Examples

-   `feat: add user login feature. refs: #42`

-   `fix: fix null pointer exception in auth. refs: #101`

-   `refactor: refactor database connection logic. refs: #88`

-   `perf: optimize batch query with bulk insert. refs: #55`

-   `test: add unit tests for wallet transfer. refs: #73`

-   `chore: bump go version to 1.23. refs: #30`

-   `mig: add email column to users table. refs: #112`

-   `ci: split unit and e2e test workflows. refs: #26`