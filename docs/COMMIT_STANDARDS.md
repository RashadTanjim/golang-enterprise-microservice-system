# Git Commit Message Standards

This repository uses a consistent commit message format so changes are easy to understand, trace, and review.

## Format

```
<type>(<scope>): <summary>

Changes:
- <what changed>

Impacted Modules:
- <service/module>

Testing:
- <tests run or N/A>
```

### Types

- `feat` - New functionality
- `fix` - Bug fix
- `refactor` - Code change with no feature/bug change
- `docs` - Documentation only
- `test` - Tests only
- `chore` - Tooling, build, or maintenance tasks

### Scopes

Use the primary area affected by the change:

- `user-service`, `order-service`, `repository-service`
- `frontend`, `common`, `monitoring`, `docs`, `infra`

## Summary Rules

- Use the imperative mood (e.g., "add", "update", "remove").
- Keep the summary under 72 characters.
- Be specific about what changed.

## Examples

```
feat(order-service): add status filter to list endpoint

Changes:
- allow filtering orders by status query parameter
- validate status values against enum

Impacted Modules:
- services/order-service/internal/handler
- services/order-service/internal/service

Testing:
- make test-order
```

```
docs(docs): add commit message standards

Changes:
- document commit message format and examples

Impacted Modules:
- docs/COMMIT_STANDARDS.md

Testing:
- N/A
```
