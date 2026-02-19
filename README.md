# GatewaySeeker

Lightning fast admin panel and hidden directory discovery tool written in Go.

## Features

- 🚀 **Multi-threaded scanning** - Configurable concurrency for maximum speed
- 📚 **Built-in wordlist** - Comprehensive list of common paths and admin panels
- 🔧 **Custom wordlist support** - Use your own wordlists
- 📝 **Extension brute-forcing** - Test multiple file extensions (.php, .asp, .bak, etc.)
- 🎯 **Status code filtering** - Hide unwanted responses (404, 403, etc.)
- 📊 **Response size analysis** - See content length for each discovery
- 🕵️ **Stealth mode** - Random delays & User-Agent rotation to avoid detection
- 🎨 **Colored console output** - Easy-to-read results with color coding
- 📁 **JSON output** - Save results for later analysis
- ⚡ **Blazing fast** - Written in Go with efficient concurrency

## Installation

```bash
# Clone the repository
git clone https://github.com/tools-tensorsecurityacademy/GatewaySeeker.git

# Go to the directory
cd GatewaySeeker

# Download dependencies
go mod download

# Build the tool
go build -o gatewayseeker main.go

# Run it
./gatewayseeker -u https://example.com