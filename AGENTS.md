# AGENTS.md

<!-- BEGIN DANNIE-GIT-IDENTITY -->
# Git Identity

When committing in this workspace, use:

```bash
git config user.name "sfzorin"
git config user.email "s.zorin@dannie.cc"
```

Commits should appear as `sfzorin <s.zorin@dannie.cc>`.
<!-- END DANNIE-GIT-IDENTITY -->

<!-- BEGIN DANNIE-GITHUB-TOKEN -->
# GitHub Token Profile

Token file: `/Users/sergeyzorin/ai/.github-tokens.env`. Never print token values and never commit that file.

Default for this folder: **gh-sfzorin** -> `GH_TOKEN_SFZORIN` for `https://github.com/sfzorin`.

When working here, load:

```bash
set -a
source /Users/sergeyzorin/ai/.github-tokens.env
set +a
export GH_TOKEN="$GH_TOKEN_SFZORIN"
```
<!-- END DANNIE-GITHUB-TOKEN -->

# Project Notes

This project is deployed to `morning30.com` from:

```text
https://github.com/sfzorin/morning30
```

Do not commit local databases, deploy keys, secrets, or generated runtime state:

```text
morning.db
morning.db-*
.codex-ssh
.env
.env.*
```
