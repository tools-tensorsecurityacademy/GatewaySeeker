package wordlist

import (
	"bufio"
	"os"
	"strings"
)

func LoadFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		
		if word != "" && !strings.HasPrefix(word, "#") {
			words = append(words, word)
		}
	}
	
	return words, scanner.Err()
}

func GetBuiltInWordlist() []string {
	return []string{
		// Admin Panels
		"admin", "administrator", "adm", "panel", "cpanel", "whm",
		"dashboard", "controlpanel", "cp", "manager", "management",
		
		// Login Pages
		"login", "signin", "auth", "authenticate", "user", "users",
		"member", "members", "account", "accounts", "profile",
		
		// Configuration
		"config", "configuration", "conf", "cfg", "settings",
		"setup", "install", "configure", "env", ".env",
		"htaccess", ".htaccess", "htpasswd", ".htpasswd",
		
		// Backups
		"backup", "bak", "temp", "tmp", "old", "new", "archive",
		"backups", "sql", "db", "database", "dump", "backup.sql",
		
		// Development
		"dev", "development", "stage", "staging", "test", "tests",
		"demo", "sandbox", "debug", "testing", "staging",
		
		// APIs
		"api", "rest", "graphql", "swagger", "swagger-ui",
		"docs", "documentation", "v1", "v2", "v3", "v4",
		
		// Content
		"uploads", "images", "img", "css", "js", "assets",
		"static", "public", "files", "media", "downloads",
		
		// CMS Specific
		"wp-admin", "wp-content", "wp-includes", "wordpress", "wp",
		"administrator", "joomla", "drupal", "magento", "shopify",
		
		// Version Control
		".git", ".svn", ".github", ".gitlab", ".idea", ".vscode",
		
		// Logs
		"logs", "log", "error_log", "access_log", "debug",
		"debug_log", "application.log", "server.log",
		
		// Sensitive
		"private", "secure", "hidden", "secret", "confidential",
		"internal", "restricted", "classified", "protected",
		
		// Shells
		"shell", "cmd", "exec", "terminal", "console", "bash",
		"shell.php", "cmd.php", "exec.php", "backdoor",
		
		// Server Info
		"server-status", "server-info", "info", "phpinfo",
		"info.php", "status", "health", "healthcheck",
		
		// Common Files
		"index", "home", "default", "main", "robots.txt",
		"sitemap.xml", "crossdomain.xml", "favicon.ico",
		
		// Web Services
		"soap", "wsdl", "xmlrpc", "rpc", "json", "ajax",
		
		// Database
		"phpmyadmin", "pma", "mysql", "dbadmin", "adminer",
		"phpPgAdmin", "pgadmin", "phpMyAdmin", "myadmin",
		
		// File Managers
		"filemanager", "files", "manager", "explorer", "browser",
		
		// Miscellaneous
		"cgi-bin", "cgi", "bin", "scripts", "includes", "classes",
		"lib", "libs", "vendor", "node_modules", "bower_components",
	}
}
