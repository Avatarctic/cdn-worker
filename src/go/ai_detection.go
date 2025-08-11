package main

import "strings"

// AI Crawlers list
var aiCrawlers = []string{
	"GPTBot",
	"OAI-SearchBot",
	"ChatGPT-User",
	"anthropic-ai",
	"ClaudeBot",
	"PerplexityBot",
	"Google-Extended",
	"Bytespider",
	"Amazonbot",
	"Applebot-Extended",
	"CCBot",
}

// IsAICrawler checks if the user agent is from an AI crawler
func IsAICrawler(userAgent string) bool {
	for _, crawler := range aiCrawlers {
		if strings.Contains(userAgent, crawler) {
			return true
		}
	}
	return false
}
