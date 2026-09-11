package main

import (
	"encoding/json"
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
		asJSON   bool
	)

	flag.StringVar(&repoPath, "C", ".", "path to the git repository")
	flag.StringVar(&since, "since", "", "only consider commits after this date (e.g. 2024-01-01)")
	flag.IntVar(&minCount, "min-count", 2, "minimum shared commits required to report a pair")
	flag.IntVar(&limit, "limit", 20, "maximum number of pairs to print")
	flag.BoolVar(&asJSON, "json", false, "print results as a JSON array instead of a table")
	flag.Parse()

	pairs, err := CoChangedPairs(repoPath, since)
	if err != nil {
		fmt.Fprintln(os.Stderr, "git-change-coupling:", err)
		os.Exit(1)
	}

	pairs = filterPairs(pairs, minCount, limit)

	if asJSON {
		if err := reportJSON(pairs); err != nil {
			fmt.Fprintln(os.Stderr, "git-change-coupling:", err)
			os.Exit(1)
		}
		return
	}
	report(pairs)
}

// filterPairs applies the min-count threshold and limit that both output
// modes share, so JSON and table output always agree on what "shown" means.
func filterPairs(pairs []Pair, minCount, limit int) []Pair {
	kept := make([]Pair, 0, limit)
	for _, p := range pairs {
		if p.Count < minCount {
			continue
		}
		kept = append(kept, p)
		if len(kept) >= limit {
			break
		}
	}
	return kept
}

func report(pairs []Pair) {
	if len(pairs) == 0 {
		fmt.Println("no file pairs met the threshold")
		return
	}
	for _, p := range pairs {
		fmt.Printf("%4d  %5.1f%%  %s <-> %s\n", p.Count, p.Strength, p.A, p.B)
	}
}

func reportJSON(pairs []Pair) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(pairs)
}
