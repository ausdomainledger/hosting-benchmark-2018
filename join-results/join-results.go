package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	pBenchTotal   = regexp.MustCompile(`Total\s+([\d\.]+)`)
	pWordPressP90 = regexp.MustCompile(`90\.000%\s+([\d\.]+m?s)`)
)

var plans []struct {
	ID       string `json:"id,omitempty"`
	Provider string `json:"provider,omitempty"`
	Name     string `json:"name,omitempty"`
	Cost     int    `json:"cost,omitempty"`
}

type results []result

func (v results) Len() int           { return len(v) }
func (v results) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v results) Less(i, j int) bool { return v[i].OverallScore < v[j].OverallScore }

type result struct {
	Provider     string      `json:"provider,omitempty"`
	Sample       string      `json:"sample,omitempty"`
	Plan         interface{} `json:"plan"`
	OverallScore int64       `json:"overall_score"`
	Results      struct {
		WP1P90           time.Duration `json:"wp_1_p90,omitempty"`
		WP5P90           time.Duration `json:"wp_5_p90,omitempty"`
		WP10P90          time.Duration `json:"wp_10_p90,omitempty"`
		DBInsertDuration time.Duration `json:"db_insert_duration,omitempty"`
		DBQueries        int           `json:"db_queries,omitempty"`
		CPUOps           int           `json:"cpu_ops,omitempty"`
		IOOpenDuration   time.Duration `json:"io_open_duration,omitempty"`
		IOSeqWrite       time.Duration `json:"io_seq_write_duration,omitempty"`
		IORandomRW       int           `json:"io_random_rw,omitempty"`
		PHP1Duration     time.Duration `json:"php1_duration,omitempty"`
		PHP2Duration     time.Duration `json:"php2_duration,omitempty"`
	} `json:"results,omitempty"`
	Env struct {
		Hostname   string      `json:"hostname,omitempty"`
		Kernel     string      `json:"kernel,omitempty"`
		CloudLinux interface{} `json:"cloudlinux,omitempty"`
		CPU        string      `json:"cpu,omitempty"`
		DBVersion  string      `json:"db_version,omitempty"`
		PHPVersion string      `json:"php_version,omitempty"`
		Server     string      `json:"server,omitempty"`
	} `json:"env,omitempty"`
}

func main() {
	var onlyMedian bool
	flag.BoolVar(&onlyMedian, "median", false, "Only output the median result for each provider")
	flag.Parse()

	ranks := results{}

	f, err := os.Open("./datasets/plans.js")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&plans); err != nil {
		panic(err)
	}

	dirs, _ := filepath.Glob("./datasets/benchmark-results/*/*")
	for _, d := range dirs {
		if stat, err := os.Stat(d); err != nil || !stat.IsDir() {
			continue
		}

		path, suffix := filepath.Split(d)
		if suffix[0] != '2' {
			continue
		}

		provider := filepath.Base(path) + "-" + suffix

		ranks = append(ranks, process(d, provider))
	}
	sort.Sort(ranks)

	if onlyMedian {
		visited := map[string][]int{}
		for i, entry := range ranks {
			visited[entry.Provider] = append(visited[entry.Provider], i)
		}
		mo := results{}
		for _, entry := range visited {
			mo = append(mo, ranks[entry[1]])
		}
		ranks = mo
		sort.Sort(ranks)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(ranks)
}

func process(dir string, provider string) result {
	name := strings.Split(provider, "-")
	r := result{Provider: name[0], Sample: name[1]}

	// bench.results
	r.Results.PHP1Duration = duration(extractRegexSingle(dir, "bench.results", pBenchTotal, true)[0] + "s")

	// database.results
	m := jsonMap(dir, "database.results", "database.php")
	r.Results.DBQueries = num(fmt.Sprintf("%v", m["queries"]))
	r.Results.DBInsertDuration = duration(fmt.Sprintf("%vs", m["insert_duration"]))

	// describe_env.results
	envBuf := globReadFile(dir, "describe_env.results")
	if err := json.Unmarshal([]byte(envBuf), &r.Env); err != nil {
		panic(err)
	}

	// micro_bench.results
	r.Results.PHP2Duration = duration(extractRegexSingle(dir, "micro_bench.results", pBenchTotal, true)[0] + "s")

	// micro_cpu.results
	m = jsonMap(dir, "micro_cpu.results", "micro_cpu.php")
	r.Results.CPUOps = num(fmt.Sprintf("%v", m["executions"]))

	// micro_io.results
	m = jsonMap(dir, "micro_io.results", "micro_io.php")
	r.Results.IOOpenDuration = duration(fmt.Sprintf("%vs", m["open_duration"]))
	r.Results.IOSeqWrite = duration(fmt.Sprintf("%vs", m["seq_write_duration"]))
	r.Results.IORandomRW = num(fmt.Sprintf("%v", m["random_rw_count"]))

	// wordpress-results-{URL}-{RATE}-{WHEN}.log
	r.Results.WP1P90 = duration(extractRegexSingle(dir, "wordpress-results-*-1-*.log", pWordPressP90, true)[0])
	r.Results.WP5P90 = duration(extractRegexSingle(dir, "wordpress-results-*-5-*.log", pWordPressP90, true)[0])
	r.Results.WP10P90 = duration(extractRegexSingle(dir, "wordpress-results-*-10-*.log", pWordPressP90, true)[0])

	r.OverallScore = int64(r.Results.WP1P90 + r.Results.WP5P90 + r.Results.WP10P90)

	for _, plan := range plans {
		if plan.ID == r.Provider {
			r.Plan = plan
			break
		}
	}
	return r
}

func globReadFile(dir, glob string) string {
	files, _ := filepath.Glob(filepath.Join(dir, glob))
	if len(files) != 1 {
		panic(fmt.Sprintf("Wrong number of files matched for glob %s+%s: %v", dir, glob, files))
	}

	buf, _ := ioutil.ReadFile(files[0])
	return string(buf)
}

func extractRegexSingle(dir, glob string, pattern *regexp.Regexp, mustMatch bool) []string {
	buf := globReadFile(dir, glob)
	if res := pattern.FindStringSubmatch(buf); len(res) < 2 && mustMatch {
		panic(fmt.Sprintf("No results for %v in %s...", pattern, buf[:256]))
	} else if len(res) >= 2 {
		return res[1:]
	} else {
		return []string{""}
	}
}

func jsonMap(dir, glob string, jsonPath string) map[string]interface{} {
	buf := globReadFile(dir, glob)
	pathPieces := strings.Split(jsonPath, "->")

	var out map[string]interface{}
	if err := json.Unmarshal([]byte(buf), &out); err != nil {
		panic(err)
	}

	for _, piece := range pathPieces {
		out = out[piece].(map[string]interface{})
	}

	return out
}

func num(in string) int {
	if v, err := strconv.Atoi(in); err != nil {
		panic(err)
	} else {
		return v
	}
}

func duration(in string) time.Duration {
	dur, err := time.ParseDuration(in)
	if err != nil {
		panic(err)
	}
	return dur
}

func sum(in []string) int {
	var out int
	for _, v := range in {
		if v == "" {
			continue
		}
		out += num(v)
	}
	return out
}
