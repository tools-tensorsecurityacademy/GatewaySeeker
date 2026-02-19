package scanner

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tools-tensorsecurityacademy/GatewaySeeker/output"
	"github.com/fatih/color"
)

type Config struct {
	TargetURL   string
	Paths       []string
	Extensions  []string
	Threads     int
	Timeout     int
	Delay       int
	StealthMode bool
	FilterCodes map[int]bool
	ShowAll     bool
}

type Result struct {
	URL         string
	StatusCode  int
	Size        int64
	SizeHuman   string
	Error       error
}

type ScanResults struct {
	TotalRequests   int
	DiscoveredCount int
	Duration        time.Duration
}

// User agents for stealth mode
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (X11; Linux i686; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_2 like Mac OS X) AppleWebKit/605.1.15",
	"Mozilla/5.0 (iPad; CPU OS 17_1_2 like Mac OS X) AppleWebKit/605.1.15",
}

func Scan(config *Config, outHandler *output.Handler) (*ScanResults, error) {
	startTime := time.Now()
	
	
	var urlsToCheck []string
	for _, path := range config.Paths {
		for _, ext := range config.Extensions {
			fullPath := path
			if ext != "" {
				fullPath = strings.TrimSuffix(path, "/") + ext
			}
			fullURL := fmt.Sprintf("%s/%s", config.TargetURL, strings.TrimLeft(fullPath, "/"))
			urlsToCheck = append(urlsToCheck, fullURL)
		}
	}

	
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:    config.Threads,
		MaxIdleConnsPerHost: config.Threads,
	}
	
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	
	jobs := make(chan string, len(urlsToCheck))
	results := make(chan Result, config.Threads)

	// Fill jobs channel
	for _, url := range urlsToCheck {
		jobs <- url
	}
	close(jobs)


	var wg sync.WaitGroup
	for i := 0; i < config.Threads; i++ {
		wg.Add(1)
		go worker(i, client, jobs, results, &wg, config)
	}

	
	go func() {
		wg.Wait()
		close(results)
	}()

	
	discoveredCount := 0
	var jsonResults []map[string]interface{}

	for res := range results {
		if res.Error != nil {
			if config.ShowAll {
				color.Red("[ERROR] %s - %v", res.URL, res.Error)
			}
			continue
		}

		
		if config.FilterCodes[res.StatusCode] && !config.ShowAll {
			continue
		}

		
		var printColor func(format string, a ...interface{})
		switch {
		case res.StatusCode >= 200 && res.StatusCode < 300:
			printColor = color.Green
			discoveredCount++
		case res.StatusCode >= 300 && res.StatusCode < 400:
			printColor = color.Yellow
		case res.StatusCode >= 400 && res.StatusCode < 500:
			printColor = color.Red
		default:
			printColor = color.Blue
		}

		
		printColor("[%d] %s - Size: %s", res.StatusCode, res.URL, res.SizeHuman)

		// Prepare for JSON output
		jsonResults = append(jsonResults, map[string]interface{}{
			"url":         res.URL,
			"status_code": res.StatusCode,
			"size_bytes":  res.Size,
			"size_human":  res.SizeHuman,
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	}

	
	duration := time.Since(startTime)

	
	if outHandler != nil && len(jsonResults) > 0 {
		err := outHandler.Save(jsonResults, ScanResults{
			TotalRequests:   len(urlsToCheck),
			DiscoveredCount: discoveredCount,
			Duration:        duration,
		})
		if err != nil {
			color.Red("[!] Failed to save JSON output: %v", err)
		}
	}

	return &ScanResults{
		TotalRequests:   len(urlsToCheck),
		DiscoveredCount: discoveredCount,
		Duration:        duration,
	}, nil
}

func worker(id int, client *http.Client, jobs <-chan string, results chan<- Result, wg *sync.WaitGroup, config *Config) {
	defer wg.Done()
	
	for url := range jobs {
		
		if config.StealthMode || config.Delay > 0 {
			delay := config.Delay
			if config.StealthMode {
				// Random delay between 100ms and 1000ms for stealth
				delay = rand.Intn(900) + 100
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			results <- Result{URL: url, Error: err}
			continue
		}

		
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Connection", "keep-alive")
		
		// Set User-Agent
		if config.StealthMode {
			ua := userAgents[rand.Intn(len(userAgents))]
			req.Header.Set("User-Agent", ua)
		} else {
			req.Header.Set("User-Agent", "GatewaySeeker/1.0")
		}

		
		resp, err := client.Do(req)
		if err != nil {
			results <- Result{URL: url, Error: err}
			continue
		}

		
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		size := int64(len(body))
		sizeHuman := formatSize(size)
		
		results <- Result{
			URL:        url,
			StatusCode: resp.StatusCode,
			Size:       size,
			SizeHuman:  sizeHuman,
			Error:      nil,
		}
	}
}

func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
}
