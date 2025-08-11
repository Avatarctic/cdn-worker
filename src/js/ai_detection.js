// AI Crawler Detection Module
// Contains logic for identifying AI crawlers

// AI Crawlers list
const aiCrawlers = [
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
];

// Check if the user agent is from an AI crawler
function isAICrawler(userAgent) {
  if (!userAgent) return false;
  return aiCrawlers.some(crawler => userAgent.includes(crawler));
}

// Export functions and constants
export { aiCrawlers, isAICrawler };
