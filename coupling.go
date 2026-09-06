package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// Pair describes how often two files were touched by the same commit.
type Pair struct {
	A, B     string
	Count    int
	Strength float64 // Count as a percentage of the less-frequently-changed file's total commits
}

// maxFilesPerCommit skips commits that touch an unusually large number of
// files (mass renames, vendoring, initial imports). Those commits inflate
// every pair count without reflecting any real relationship between files.
const maxFilesPerCommit = 50

// CoChangedPairs walks the commit history of repoPath and counts, for every
// pair of files, how many commits touched both of them.
func CoChangedPairs(repoPath, since string) ([]Pair, error) {
	args := []string{"-C", repoPath, "log", "-z", "--pretty=format:%H", "--name-only"}
	if since != "" {
		args = append(args, "--since="+since)
	}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("running git log: %w", err)
	}

	fileCommits := map[string]int{}
	pairCommits := map[[2]string]int{}

	// git log -z separates whole commit records with NUL instead of the
	// usual blank-line spacing, so each record here is "hash\nfile\nfile...".
	for _, rec := range strings.Split(string(out), "\x00") {
		rec = strings.TrimSpace(rec)
		if rec == "" {
			continue
		}
		lines := strings.Split(rec, "\n")

		files := make([]string, 0, len(lines)-1)
		for _, line := range lines[1:] {
			if line = strings.TrimSpace(line); line != "" {
				files = append(files, line)
			}
		}
		if len(files) < 2 || len(files) > maxFilesPerCommit {
			continue
		}
		sort.Strings(files)

		for _, f := range files {
			fileCommits[f]++
		}
		for i := 0; i < len(files); i++ {
			for j := i + 1; j < len(files); j++ {
				pairCommits[[2]string{files[i], files[j]}]++
			}
		}
	}

	pairs := make([]Pair, 0, len(pairCommits))
	for key, count := range pairCommits {
		a, b := key[0], key[1]
		smaller := fileCommits[a]
		if fileCommits[b] < smaller {
			smaller = fileCommits[b]
		}
		pairs = append(pairs, Pair{
			A:        a,
			B:        b,
			Count:    count,
			Strength: float64(count) / float64(smaller) * 100,
		})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Strength > pairs[j].Strength
	})

	return pairs, nil
}
