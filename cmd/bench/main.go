package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type bench_config struct {
	base_url    string
	concurrency int
	duration    time.Duration
	identifier  string
	password    string
}

type result struct {
	latency time.Duration
	ok      bool
	code    int
	path    string
}

func main() {
	cfg := parse_flags()
	fmt.Printf("[bench] base=%s concurrency=%d duration=%s\n", cfg.base_url, cfg.concurrency, cfg.duration)

	client := &http.Client{ Timeout: 15 * time.Second }

	// pre-login tokens per worker
	tokens := make([]string, cfg.concurrency)
	for i := 0; i < cfg.concurrency; i++ {
		acc, err := do_login(client, cfg.base_url, cfg.identifier, cfg.password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
			os.Exit(1)
		}
		tokens[i] = acc
	}

	var wg sync.WaitGroup
	res_ch := make(chan result, 10000)
	var stop int32
	end := time.Now().Add(cfg.duration)

	paths := []string{"/api/v1/users/me", "/api/v1/users?limit=20&offset=0"}

	for i := 0; i < cfg.concurrency; i++ {
		wg.Add(1)
		acc := tokens[i]
		go func(idx int) {
			defer wg.Done()
			for atomic.LoadInt32(&stop) == 0 && time.Now().Before(end) {
				for _, p := range paths {
					start := time.Now()
					code := do_get_auth(client, cfg.base_url+p, acc)
					ok := code >= 200 && code < 300
					res_ch <- result{latency: time.Since(start), ok: ok, code: code, path: p}
				}
			}
		}(i)
	}

	// collector
	go func() {
		wg.Wait()
		close(res_ch)
	}()

	stats := collect(res_ch)
	atomic.StoreInt32(&stop, 1)

	print_stats(stats)
}

func parse_flags() bench_config {
	base := flag.String("base", "http://localhost:8080", "base url of server")
	conc := flag.Int("concurrency", 10, "number of concurrent workers")
	dur := flag.Duration("duration", 10*time.Second, "total test duration")
	id := flag.String("identifier", "admin@example.com", "login identifier (email or username)")
	pwd := flag.String("password", "admin123", "login password")
	flag.Parse()
	return bench_config{ base_url: strings.TrimRight(*base, "/"), concurrency: *conc, duration: *dur, identifier: *id, password: *pwd }
}

func do_login(c *http.Client, base, identifier, password string) (string, error) {
	body := map[string]string{"identifier": identifier, "password": password}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", base+"/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		bb, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login status=%d body=%s", resp.StatusCode, string(bb))
	}
	var out struct{ AccessToken string `json:"access_token"` }
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
	if out.AccessToken == "" { return "", fmt.Errorf("empty access_token") }
	return out.AccessToken, nil
}

func do_get_auth(c *http.Client, url, token string) int {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.Do(req)
	if err != nil { return 0 }
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

type stats_summary struct {
	total   int
	errors  int
	codes   map[int]int
	p50     time.Duration
	p95     time.Duration
	p99     time.Duration
	avg     time.Duration
	max     time.Duration
	min     time.Duration
	per_path map[string]int
}

func collect(ch <-chan result) stats_summary {
	var lats []time.Duration
	codes := make(map[int]int)
	per_path := make(map[string]int)
	var total, errors int
	var min time.Duration = time.Hour
	var max time.Duration
	var sum time.Duration

	for r := range ch {
		total++
		if !r.ok { errors++ }
		codes[r.code]++
		per_path[r.path]++
		lats = append(lats, r.latency)
		if r.latency < min { min = r.latency }
		if r.latency > max { max = r.latency }
		sum += r.latency
	}

	if total == 0 { return stats_summary{} }
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })
	p := func(q float64) time.Duration {
		if len(lats) == 0 { return 0 }
		idx := int(math.Ceil(q*float64(len(lats))) - 1)
		if idx < 0 { idx = 0 }
		if idx >= len(lats) { idx = len(lats)-1 }
		return lats[idx]
	}
	avg := time.Duration(int64(sum) / int64(total))
	return stats_summary{
		total: total,
		errors: errors,
		codes: codes,
		p50: p(0.50),
		p95: p(0.95),
		p99: p(0.99),
		avg: avg,
		max: max,
		min: min,
		per_path: per_path,
	}
}

func print_stats(s stats_summary) {
	fmt.Println("\n[bench] results")
	fmt.Printf("  total: %d\n", s.total)
	fmt.Printf("  errors: %d\n", s.errors)
	fmt.Printf("  latency: min=%s avg=%s p50=%s p95=%s p99=%s max=%s\n", s.min, s.avg, s.p50, s.p95, s.p99, s.max)
	fmt.Printf("  status codes:\n")
	// stable order
	keys := make([]int, 0, len(s.codes))
	for k := range s.codes { keys = append(keys, k) }
	sort.Ints(keys)
	for _, k := range keys { fmt.Printf("    %d: %d\n", k, s.codes[k]) }
	fmt.Printf("  per path:\n")
	paths := make([]string, 0, len(s.per_path))
	for p := range s.per_path { paths = append(paths, p) }
	sort.Strings(paths)
	for _, p := range paths { fmt.Printf("    %s: %d\n", p, s.per_path[p]) }
}
