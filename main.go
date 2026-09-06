package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		since    string
		limit    int
		minCount int
		repoPath string
	)

	flag.StringVar(&repoPath, "C", ".", "path to the git repository")
	flag.StringVar(&since, "since", "", "only consider commits after this date (e.g. 2024-01-01)")
	flag.IntVar(&minCount, "min-count", 2, "minimum shared commits required to report a pair")
	flag.IntVar(&limit, "limit", 20, "maximum number of pairs to print")
	flag.Parse()

	pairs, err := CoChangedPairs(repoPath, since)
	if err != nil {
		fmt.Fprintln(os.Stderr, "git-change-coupling:", err)
		os.Exit(1)
	}

	report(pairs, minCount, limit)
}

func report(pairs []Pair, minCount, limit int) {
	shown := 0
	for _, p := range pairs {
		if p.Count < minCount {
			continue
		}
		fmt.Printf("%4d  %5.1f%%  %s <-> %s\n", p.Count, p.Strength, p.A, p.B)
		shown++
		if shown >= limit {
			break
		}
	}
	if shown == 0 {
		fmt.Println("no file pairs met the threshold")
	}
}
