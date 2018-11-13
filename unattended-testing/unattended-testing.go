package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type hostDef struct {
	BaseURL     string
	BenchSuffix string
	WPSuffix    string
}

type testDef struct {
	Name       string
	SleepAfter time.Duration
}

var (
	benchmarks = map[string]hostDef{
		"ventraip":       {"http://blog.id-rsa.pub", "/", "/wordpress/"},
		"micron21":       {"http://blug.id-rsa.pub", "/.bench/", "/wordpress/"},
		"zuver":          {"http://bleg.id-rsa.pub", "/.bench/", "/wordpress/"},
		"netorigin":      {"http://blig.id-rsa.pub", "/.bench/", "/wordpress/"},
		"melbourneit":    {"http://mapp.id-rsa.pub", "/.bench/", "/wordpress/"},
		"digitalpacific": {"http://zapp.id-rsa.pub", "/.bench/", "/wordpress/"},
		"netvirtue":      {"http://blyg.id-rsa.pub", "/.bench/", "/wordpress/"},
		"panthur":        {"http://blag.id-rsa.pub", "/", "/wordpress/"},
	}
	tests = []testDef{
		{"describe_env", 5 * time.Second},
		{"bench", 10 * time.Second},
		{"micro_bench", 30 * time.Second},
		{"micro_cpu", 2 * time.Minute},
		{"micro_io", 2 * time.Minute},
		{"database", 2 * time.Minute},
	}
	wg       sync.WaitGroup
	cl       *http.Client
	wrkPath  *string
	wrkRates = []int{1, 5, 10}
	wrkMu    sync.Mutex // only one instance of wrk2 may run at once
	wrkSleep = 5 * time.Minute
)

func main() {
	onlyOpt := flag.String("only", "", "CSV of host names to run")
	intervalHours := flag.Int("interval", 6, "Number of hours between benchmarks")
	wrkPath = flag.String("wrk-bin", "./wrk2/wrk", "Path to the wrk2 binary")
	flag.Parse()

	if *onlyOpt != "" {
		onlyList := strings.Split(*onlyOpt, ",")
		for k := range benchmarks {
			found := false
			for _, n := range onlyList {
				if k == n {
					found = true
					break
				}
			}
			if !found {
				log.Printf("Skipping %s due to -only", k)
				delete(benchmarks, k)
			}
		}
	}

	cl = &http.Client{}

	for {
		log.Printf("Starting set")

		wg.Add(len(benchmarks))
		for n, p := range benchmarks {
			go runBenchmark(n, p)
		}

		wg.Wait()

		log.Printf("Completed set, will sleep %d hours", *intervalHours)
		time.Sleep(time.Duration(*intervalHours) * time.Hour)
	}
}

func runBenchmark(name string, params hostDef) {
	defer wg.Done()

	ts := time.Now().Format("200601021504")

	resultsDir := filepath.Join("./results/", name, ts)
	if err := os.MkdirAll(resultsDir, 0700); err != nil {
		log.Printf("[ERROR] Failed to create results dir: %v", err)
		return
	}

	// Run micro-benchmarks
	for n, test := range tests {
		u := params.BaseURL + params.BenchSuffix + test.Name + ".php"
		if err := fetchAndPersist(u, filepath.Join(resultsDir, test.Name+".results")); err != nil {
			log.Printf("[ERROR] Failed to fetch and save %s: %v", u, err)
		} else {
			log.Printf("[DEBUG] Completed %s for %s, sleeping %v", test.Name, name, test.SleepAfter)
		}
		if n < len(tests)-1 {
			time.Sleep(test.SleepAfter)
		}
	}

	// Run wrk2
	u := params.BaseURL + params.WPSuffix
	for n, rate := range wrkRates {
		if err := runWrkAndPersist(
			u, rate, filepath.Join(resultsDir, fmt.Sprintf("wordpress-results-%s-%d-%s.log", name, rate, ts))); err != nil {
			log.Printf("[ERROR] Failed to run wrk2 on %s @ %d/sec: %v", u, rate, err)
		} else {
			log.Printf("[DEBUG] Completed wrk2 for %s @ %d/sec, will sleep 5m", u, rate)
		}
		if n < len(wrkRates)-1 {
			time.Sleep(wrkSleep)
		}
	}

}

func fetchAndPersist(u, resultsFile string) error {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	resp, err := cl.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("danger: %s is a 404", u)
	}

	f, err := os.OpenFile(resultsFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	r := io.LimitReader(resp.Body, 8*1024*1024)
	if _, err := io.Copy(f, r); err != io.EOF {
		return err
	}

	return nil
}

func runWrkAndPersist(u string, rate int, resultsFile string) error {
	wrkMu.Lock()
	defer wrkMu.Unlock()

	f, err := os.OpenFile(resultsFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c := exec.CommandContext(ctx, *wrkPath,
		"-H", "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:64.0) Gecko/20100101 Firefox/64.0",
		"--latency", "-d60s", fmt.Sprintf("-R%d", rate), "-t", "1", "-c", strconv.Itoa(rate), u)
	c.Stdout = f
	c.Stderr = f

	return c.Run()
}
