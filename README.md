# git-change-coupling

Import graphs and call graphs only show the dependencies a compiler can see.
They miss the ones that live in people's heads: the handler that always needs
its test updated, the two YAML files that have to stay in sync, the migration
that never ships without a matching model change. Those files aren't linked
by any `import`, but they show up together in commit after commit.

`git-change-coupling` mines `git log` for exactly that pattern. For every
pair of files, it counts how many commits touched both, and reports the
pairs that change together most often. High coupling with no visible code
relationship is usually worth investigating - it's either a missing
abstraction or an implicit contract that isn't written down anywhere.

## Usage

```
$ git-change-coupling -C ~/code/myrepo -min-count 3 -limit 10
  12  85.7%  internal/api/handler.go <-> internal/api/handler_test.go
   9  60.0%  config/prod.yaml <-> config/staging.yaml
   7  43.8%  db/schema.sql <-> internal/models/user.go
   5  71.4%  cmd/server/main.go <-> cmd/server/flags.go
```

Each line shows:

- the number of commits that touched both files
- that count as a percentage of the less-active file's total commits (a rough
  measure of how "coupled" the pair is, not just how busy the files are)
- the two file paths

## Flags

| Flag         | Default | Meaning                                              |
|--------------|---------|-------------------------------------------------------|
| `-C`         | `.`     | path to the git repository                            |
| `-since`     | (none)  | only consider commits after this date, e.g. `2024-01-01` |
| `-min-count` | `2`     | minimum shared commits required to report a pair       |
| `-limit`     | `20`    | maximum number of pairs to print                       |

Commits that touch more than 50 files (mass renames, initial imports,
vendoring) are skipped, since they inflate every pair count without meaning
anything.

## Building

```
go build -o git-change-coupling .
```

No third-party dependencies - standard library only.

## Limitations

This only sees what got committed together, not why. Two files can be
coupled because they belong together, or because someone forgot to split an
unrelated change into two commits. Read the results, don't just automate on
top of them.
