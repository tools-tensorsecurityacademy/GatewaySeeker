package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tools-tensorsecurityacademy/GatewaySeeker/output"
	"github.com/tools-tensorsecurityacademy/GatewaySeeker/scanner"
	"github.com/tools-tensorsecurityacademy/GatewaySeeker/wordlist"
	"github.com/fatih/color"
)

var (
	version = "1.0.0"
	commit  = "none"
	date    = "unknown"
)

func main() {
	
	color.Cyan(`
   ╔══════════════════════════════════════════════════════════╗
   ║     GatewaySeeker - Lightning Fast Directory Discovery   ║
   ║                    Version ` + version + `                           ║
   ╚══════════════════════════════════════════════════════════╝
	`)

	
	targetURL := flag.String("u", "", "Target URL (required)")
	wordlistPath := flag.String("w", "", "Path to custom wordlist file")
	threads := flag.Int("t", 20, "Number of concurrent threads")
	extensions := flag.String("ext", "php,asp,aspx,jsp,do,bak,txt,html", "File extensions to try (comma separated)")
	filterCodes := flag.String("fc", "", "Hide status codes (example: 404,403)")
	delay := flag.Int("delay", 0, "Delay between requests in milliseconds")
	outputFile := flag.String("o", "", "Save results to JSON file")
	stealth := flag.Bool("stealth", false, "Enable stealth mode (random delays + user agents)")
	timeout := flag.Int("timeout", 10, "HTTP request timeout in seconds")
	showAll := flag.Bool("all", false, "Show all responses including errors")
	showVersion := flag.Bool("version", false, "Show version information")
	
	flag.Parse()

	
	if *showVersion {
		fmt.Printf("GatewaySeeker version %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
		os.Exit(0)
	}

	
	if *targetURL == "" {
		color.Red("[!] Error: Target URL is required")
		color.Yellow("[*] Usage: gatewayseeker -u https://example.com")
		color.Yellow("[*] For help: gatewayseeker -h")
		os.Exit(1)
	}

	
	targetURLStr := *targetURL
	if !strings.HasPrefix(targetURLStr, "http://") && !strings.HasPrefix(targetURLStr, "https://") {
		targetURLStr = "https://" + targetURLStr
		color.Yellow("[*] No protocol specified, using HTTPS: %s", targetURLStr)
	}
	targetURLStr = strings.TrimRight(targetURLStr, "/")

	
	extList := []string{""} // Always include no extension
	if *extensions != "" {
		for _, ext := range strings.Split(*extensions, ",") {
			ext = strings.TrimSpace(ext)
			if ext != "" {
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				extList = append(extList, ext)
			}
		}
	}

	
	filterCodeMap := make(map[int]bool)
	if *filterCodes != "" {
		for _, code := range strings.Split(*filterCodes, ",") {
			var c int
			fmt.Sscanf(strings.TrimSpace(code), "%d", &c)
			if c > 0 {
				filterCodeMap[c] = true
			}
		}
	}

	
	var paths []string
	var err error
	
	if *wordlistPath != "" {
		color.Blue("[*] Loading wordlist from: %s", *wordlistPath)
		paths, err = wordlist.LoadFromFile(*wordlistPath)
		if err != nil {
			color.Red("[!] Failed to load wordlist: %v", err)
			os.Exit(1)
		}
	} else {
		color.Blue("[*] Using built-in wordlist")
		paths = wordlist.GetBuiltInWordlist()
	}

	
	color.Cyan("\n[*] Scan Configuration:")
	color.White("    Target: %s", targetURLStr)
	color.White("    Words to check: %d", len(paths))
	color.White("    Total URLs: %d (with extensions)", len(paths)*len(extList))
	color.White("    Threads: %d", *threads)
	color.White("    Extensions: %s", strings.Join(extList[1:], ", "))
	
	if *stealth {
		color.Yellow("    Stealth mode: ON")
	}
	if *delay > 0 {
		color.Yellow("    Delay: %dms", *delay)
	}
	
	color.Cyan("\n[*] Starting scan...\n")

	
	config := &scanner.Config{
		TargetURL:    targetURLStr,
		Paths:        paths,
		Extensions:   extList,
		Threads:      *threads,
		Timeout:      *timeout,
		Delay:        *delay,
		StealthMode:  *stealth,
		FilterCodes:  filterCodeMap,
		ShowAll:      *showAll,
	}

	
	outHandler := output.NewHandler(*outputFile)

	
	results, err := scanner.Scan(config, outHandler)
	if err != nil {
		color.Red("[!] Scan failed: %v", err)
		os.Exit(1)
	}

	
	color.Cyan("\n[*] Scan Complete!")
	color.White("    Total URLs checked: %d", results.TotalRequests)
	color.White("    Discovered: %d", results.DiscoveredCount)
	color.White("    Time taken: %.2f seconds", results.Duration.Seconds())
	color.White("    Requests/sec: %.2f", float64(results.TotalRequests)/results.Duration.Seconds())

	if *outputFile != "" {
		color.Green("[+] Results saved to: %s", *outputFile)
	}
}