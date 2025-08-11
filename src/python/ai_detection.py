# AI Crawler Detection Module
# Contains logic for identifying AI crawlers

# AI Crawlers list
AI_CRAWLERS = [
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
]

def is_ai_crawler(user_agent: str) -> bool:
    """Check if the user agent is from an AI crawler"""
    return any(crawler in user_agent for crawler in AI_CRAWLERS)
