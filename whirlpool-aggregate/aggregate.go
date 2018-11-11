package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

var (
	patterns []*regexp.Regexp
	threads  map[string]thread
)

type thread struct {
	Title string   `json:"title"`
	Posts [][]post `json:"posts"`
}

type post struct {
	Text string `json:"text"`
	Date string `json:"date"`
}

func main() {
	if err := readPatterns(); err != nil {
		panic(err)
	}

	f, err := os.Open("./datasets/whirlpool-threads.js")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&threads); err != nil {
		panic(err)
	}

	counts := map[string]int{}

	for _, thread := range threads {
		for _, post := range thread.Posts[0] {
			for _, pattern := range patterns {
				if pattern.MatchString(strings.ToLower(post.Text)) {
					counts[pattern.String()] = counts[pattern.String()] + 1
				}
			}
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	for p, count := range counts {
		fmt.Fprintf(w, "%s\t%d\n", p, count)
	}

}

func readPatterns() error {
	f, err := os.Open("./datasets/company-name-patterns")
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if t := sc.Text(); t[0] != '#' {
			patterns = append(patterns, regexp.MustCompile(t))
		}
	}
	return sc.Err()
}
