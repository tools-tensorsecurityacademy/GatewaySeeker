# GatewaySeeker 🔍

**Lightning-fast admin panel & hidden directory discovery tool** written in Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/tools-tensorsecurityacademy/GatewaySeeker)](https://goreportcard.com/report/github.com/tools-tensorsecurityacademy/GatewaySeeker)

Ever wasted hours manually guessing `/admin`, `/login`, `/backup`, etc.?  
Yeah… we’ve all been there.

GatewaySeeker is your fast, stealthy, and colorful companion that does the boring work for you — so you can focus on the interesting findings.

##  Why people like it

-  **Very fast** — thanks to Go + concurrent goroutines (hundreds of req/s easily)
-  **Beautiful colored output** — green = jackpot, red = nope, yellow = maybe…
-  **Stealth mode** — random delays + rotating User-Agents
-  **Works out of the box** — comes with 150+ common admin/hidden paths
-  **Super customizable** — your wordlists, your extensions, your filters

## Features at a glance

| Feature                | What it does                                      | Built-in? |
|------------------------|---------------------------------------------------|-----------|
| Multi-threading        | Scale to 10, 50, 200+ threads                     | Yes       |
| Extension brute-forcing| Try `.php`, `.bak`, `.old`, `.env`, etc.          | Yes       |
| Status code filtering  | Hide 404s, 403s, …                                | Yes       |
| Response size display  | Know if it's a real page or tiny error            | Yes       |
| Stealth mode           | Random delays + User-Agent rotation               | Yes       |
| JSON output            | Easy to save / parse results                      | Yes       |
| Built-in wordlist      | 150+ realistic admin & debug paths                | Yes       |
| Custom wordlists       | Feed it whatever you want                         | Yes       |

##  Installation

### Easiest way (recommended)

```bash
go install github.com/tools-tensorsecurityacademy/GatewaySeeker@latest
```
Make sure $GOPATH/bin is in your $PATH:
```
Bashexport PATH=$$   PATH:   $$(go env GOPATH)/bin
```
Now just run:
```
BashGatewaySeeker -u https://example.com
```

### Build from source
```
git clone https://github.com/tools-tensorsecurityacademy/GatewaySeeker.git
cd GatewaySeeker

go mod tidy
go build -o gatewayseeker

# Optional: put it somewhere nice
sudo mv gatewayseeker /usr/local/bin/
```

### Quick try without installing:
```
git clone https://github.com/tools-tensorsecurityacademy/GatewaySeeker.git

cd GatewaySeeker

go run . -u https://example.com
```

### Quick Start Examples
#### Basic scan
```
 GatewaySeeker -u https://target.com
```
Go fast
```
GatewaySeeker -u https://target.com -t 80
```
Sneaky & polite
```
GatewaySeeker -u https://target.com -stealth -delay 400 -t 15
```
Use your own wordlist + extensions
```
GatewaySeeker -u https://target.com -w mylist.txt -ext php,asp,aspx,bak,old,env
```
Hide noise & save results
```
GatewaySeeker -u https://target.com -fc 404,403,429 -o results.json
```

## All Command-line Flags
| **Flag**   | **Description**                     | **Default**                            |
| ---------- | ----------------------------------- | -------------------------------------- |
| `-u`       | **Target URL** (required)           | —                                      |
| `-w`       | Custom wordlist file                | built-in (~150 paths)                  |
| `-t`       | Number of concurrent threads        | `20`                                   |
| `-ext`     | Extensions to try (comma-separated) | `php,asp,aspx,jsp,do,bak,txt,html,old` |
| `-fc`      | Filter out these status codes       | —                                      |
| `-delay`   | Milliseconds delay between requests | `0`                                    |
| `-stealth` | Enable random delays + UA rotation  | `false`                                |
| `-timeout` | HTTP timeout (seconds)              | `10`                                   |
| `-o`       | Save results to JSON file           | —                                      |
| `-all`     | Show failed / error requests too    | `false`                                |
| `-version` | Show version and exit               | —                                      |

## What the output looks like
```
[200] https://example.com/admin       Size: 14.2 KB    ← jackpot!
[301] https://example.com/backup      Size: 0 B        ← interesting redirect
[403] https://example.com/config.php  Size: 1.1 KB     ← forbidden, but exists 👀
[404] https://example.com/test        Size: 0 B        (filtered by default)
```

## Pro Tips
- Start with `-t 20–30`, then increase if the server is friendly
- Use `-stealth` + `-delay 300–800` when facing WAFs / rate-limiting
- Always filter 404s with `-fc 404` (sometimes also 403/429)
- Targeted wordlists > generic ones — bring your own!
- Watch 301/302 redirects — they often hide the real treasure

## Responsible Use
This is a powerful tool — only use it on targets you own or have explicit written permission to test.
Unauthorized scanning is illegal in most jurisdictions.

## Contributing
Love to have you!
1. Fork it
2. Create your branch (git checkout -b feature/amazing-thing)
3. Commit (git commit -m 'Add amazing thing')
4. Push (git push origin feature/amazing-thing)
5. Open a Pull Request

Popular ideas right now:
- Add support for custom headers
- SOCKS5 / HTTP proxy support
- More built-in wordlists (especially framework-specific)
- Better rate-limit detection & auto-throttle

## License
MIT License — do cool stuff, just don’t blame me if something breaks 
Happy (ethical) hunting! 
