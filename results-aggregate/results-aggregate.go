package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	durationErrorSentinel = time.Duration(math.MaxInt64)
	countErrorSentinel    = math.MaxUint64
)

var (
	results       = map[string]*result{}
	pBenchTotal   = regexp.MustCompile(`Total\s+([\d\.]+)`)
	pWordPressP90 = regexp.MustCompile(`90\.000%\s+([\d\.]+m?s?)`)
)

type result struct {
	Provider     string                 `json:"provider,omitempty"`
	Plan         interface{}            `json:"plan,omitempty"`
	Measurements map[string]interface{} `json:"measurements,omitempty"`
	Env          struct {
		Hostname   string      `json:"hostname,omitempty"`
		Kernel     string      `json:"kernel,omitempty"`
		CloudLinux interface{} `json:"cloudlinux,omitempty"`
		CPU        string      `json:"cpu,omitempty"`
		DBVersion  string      `json:"db_version,omitempty"`
		PHPVersion string      `json:"php_version,omitempty"`
		Server     string      `json:"server,omitempty"`
	} `json:"env,omitempty"`
}

func (r *result) addTime(measureName string, d time.Duration) {
	v, ok := r.Measurements[measureName].(timeMeasurement)
	if !ok {
		v = timeMeasurement{}
	}

	if d == durationErrorSentinel {
		v.Errors++
	} else {
		v.R = append(v.R, d)
	}

	sort.Sort(v)

	v.Median = v.R[len(v.R)/2]

	r.Measurements[measureName] = v
}

func (r *result) addCounter(measureName string, c uint64) {
	v, ok := r.Measurements[measureName].(counterMeasurement)
	if !ok {
		v = counterMeasurement{}
	}

	if c == countErrorSentinel {
		v.Errors++
	} else {
		v.R = append(v.R, c)
	}

	sort.Sort(v)

	v.Median = v.R[len(v.R)/2]

	r.Measurements[measureName] = v
}

type timeMeasurement struct {
	R      []time.Duration
	Median time.Duration
	Errors int
}
type counterMeasurement struct {
	R      []uint64
	Median uint64
	Errors int
}

func (m timeMeasurement) Len() int           { return len(m.R) }
func (m timeMeasurement) Less(i, j int) bool { return m.R[i] < m.R[j] }
func (m timeMeasurement) Swap(i, j int)      { m.R[i], m.R[j] = m.R[j], m.R[i] }

func (m counterMeasurement) Len() int           { return len(m.R) }
func (m counterMeasurement) Less(i, j int) bool { return m.R[i] < m.R[j] }
func (m counterMeasurement) Swap(i, j int)      { m.R[i], m.R[j] = m.R[j], m.R[i] }

func main() {
	matches, err := filepath.Glob("./datasets/benchmark-results/*/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, match := range matches {
		stat, err := os.Stat(match)
		if err != nil {
			log.Fatal(err)
		}

		if !stat.IsDir() || stat.Name()[0] != '2' {
			continue
		}

		process(match)
	}

	var plans []struct {
		ID       string `json:"id,omitempty"`
		Provider string `json:"provider,omitempty"`
		Name     string `json:"name,omitempty"`
		Cost     int    `json:"cost,omitempty"`
	}
	f, err := os.Open("./datasets/plans.js")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&plans); err != nil {
		log.Fatal(err)
	}
	for _, plan := range plans {
		results[plan.ID].Plan = plan
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	e.Encode(results)
}

func process(resultPath string) {
	splitPath := strings.Split(resultPath, string(filepath.Separator))
	host := splitPath[2]

	r := getHost(host)

	envBuf := fread(resultPath, "describe_env.results")
	if err := json.Unmarshal([]byte(envBuf), &r.Env); err != nil {
		log.Fatal(err)
	}

	r.addTime("PHP1Duration", duration(regex(resultPath, "bench.results", pBenchTotal)+"s"))
	r.addTime("PHP2Duration", duration(regex(resultPath, "micro_bench.results", pBenchTotal)+"s"))
	r.addTime("WP1P90", duration(regex(resultPath, "wordpress-results-*-1-*.log", pWordPressP90)))
	r.addTime("WP5P90", duration(regex(resultPath, "wordpress-results-*-5-*.log", pWordPressP90)))
	r.addTime("WP10P90", duration(regex(resultPath, "wordpress-results-*-10-*.log", pWordPressP90)))

	j := readJSON(resultPath, "database.results", "database.php")
	r.addCounter("DBQueries", count(fmt.Sprintf("%v", j["queries"])))
	r.addTime("DBInsertDuration", duration(fmt.Sprintf("%vs", j["insert_duration"])))

	j = readJSON(resultPath, "micro_io.results", "micro_io.php")
	r.addTime("IOOpenDuration", duration(fmt.Sprintf("%vs", j["open_duration"])))
	r.addTime("IOSeqWrite", duration(fmt.Sprintf("%vs", j["seq_write_duration"])))
	r.addCounter("IORandomRW", count(fmt.Sprintf("%v", j["random_rw_count"])))

	j = readJSON(resultPath, "micro_cpu.results", "micro_cpu.php")
	r.addCounter("CPUOps", count(fmt.Sprintf("%v", j["executions"])))

}

func getHost(name string) *result {
	if v, ok := results[name]; ok {
		return v
	}
	v := &result{Provider: name, Measurements: map[string]interface{}{}}
	results[name] = v
	return v
}

func fread(dir, glob string) string {
	matches, err := filepath.Glob(filepath.Join(dir, glob))
	if err != nil {
		panic(err)
	}

	buf, err := ioutil.ReadFile(matches[0])
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func regex(dir, path string, p *regexp.Regexp) string {
	s := fread(dir, path)
	res := p.FindStringSubmatch(s)
	if len(res) != 2 {
		log.Fatal("bad regex", p, s)
	}
	return res[1]
}

func duration(s string) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("ERR duration(%s) = %v", s, err)
		return durationErrorSentinel
	}
	return dur
}

func count(s string) uint64 {
	c, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return countErrorSentinel
	}
	return c
}

func readJSON(dir, path, topKey string) map[string]interface{} {
	out := map[string]interface{}{}
	s := fread(dir, path)
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return out
	}
	return out[topKey].(map[string]interface{})
}
