# Contributing

## Branches

All work happens on short-lived branches merged into `main`.

| Prefix | Use for |
| --- | --- |
| `feat/` | new functionality |
| `fix/` | bug fixes |
| `chore/` | dependencies, CI, docs, tooling |

Format: `<prefix>/<kebab-case-description>` — e.g. `feat/shopping-list-export`.

## Commits

[Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
```

`type` is one of `feat`, `fix`, `chore`. Append `!` after the scope for breaking changes.

```
feat(pantry): cycle item status on tap
fix(api): reject slugs longer than 60 characters
chore(deps): bump Angular to 20.3
feat(api)!: rename items endpoint
```

The changelog is generated from these commits with [git-cliff](https://git-cliff.org/):

```bash
make -C backend changelog
```

## Before opening a pull request

```bash
make lint
make test
```

Pull requests are squash-merged: one PR produces exactly one commit on `main`.
