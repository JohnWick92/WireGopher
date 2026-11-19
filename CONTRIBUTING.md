# Contributing to WireGopher

## Code of Conduct

Be respectful, constructive, and inclusive. We're all here to build something useful together. Harassment, discrimination, or hostile behavior of any kind will not be tolerated.

## Code Style & Formatting

We enforce formatting and linting with [`pre-commit`](https://pre-commit.com/). **This is mandatory** — the CI runs the same checks on every PR, so unformatted code cannot be merged.

### One-time setup

```bash
pip install pre-commit
pre-commit install
```

After this, the hook runs **automatically on every `git commit`** — no extra command needed. It will format your Go files with `gofmt` and block the commit if anything was changed, so you can review and re-stage the fix.

### Rules

- **All Go code must be `gofmt`-clean.** No exceptions.
- Run `pre-commit run --all-files` manually before pushing to catch anything the hook missed (e.g., files added but not yet committed).
- `git commit --no-verify` is discouraged — the CI will catch it and your PR will fail.
