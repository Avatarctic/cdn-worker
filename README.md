# Cloudflare Worker for Detection of AI Crawlers

A sophisticated Cloudflare Worker implementation available in three programming languages (Go, JavaScript, Python) that provides AI crawler detection, transparent proxying, and comprehensive logging capabilities.

## 🚀 Features

### Core Functionality
- **AI Crawler Detection**: Automatically identifies and serves optimized content to AI crawlers (GPTBot, ClaudeBot, etc.)
- **Transparent Proxying**: Forwards regular traffic to origin server seamlessly
- **Comprehensive Logging**: Real-time request logging to external service
- **Multi-Language Support**: Identical functionality across Go, JavaScript, and Python

### AI Crawler Optimization
- Serves lightweight, AI-optimized HTML content
- Reduces bandwidth usage for AI training requests
- Maintains full SEO compatibility

### Future-Ready Architecture
- **Metrics & Analytics**: Preserved in `extra/` folders for future activation
- **Heuristic Detection**: Advanced bot detection beyond user-agent matching
- **Content Processing**: Strategy and Template patterns for dynamic content optimization

## 📁 Project Structure

```
cloudflare-worker/
├── README.md                 # This file
├── wrangler.toml            # Cloudflare Workers configuration
└── src/
    ├── go/                  # Go implementation
    │   ├── main.go          # Main worker logic
    │   ├── ai_detection.go  # AI crawler detection
    │   ├── logging.go       # Request logging
    │   ├── README.md        # Go-specific documentation
    │   └── extra/           # Future functionality
    │       ├── metrics.go               # Analytics & metrics
    │       ├── heuristic_detection.go   # Advanced bot detection
    │       └── content_processor.go     # Content optimization
    ├── js/                  # JavaScript implementation (Production)
    │   ├── main.js          # Main worker logic
    │   ├── ai_detection.js  # AI crawler detection
    │   ├── logging.js       # Request logging
    │   ├── README.md        # JavaScript-specific documentation
    │   └── extra/           # Future functionality
    │       ├── metrics.js               # Analytics & metrics
    │       ├── heuristic_detection.js   # Advanced bot detection
    │       └── content_processor.js     # Content optimization
    └── python/              # Python implementation
        ├── main.py          # Main worker logic
        ├── ai_detection.py  # AI crawler detection
        ├── logging.py       # Request logging
        ├── README.md        # Python-specific documentation
        └── extra/           # Future functionality
            ├── metrics.py               # Analytics & metrics
            ├── heuristic_detection.py   # Advanced bot detection
            └── content_processor.py     # Content optimization

```

## 🛠 Quick Start

### Deploy to Cloudflare Workers

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd cloudflare-worker
   ```

2. **Configure Wrangler**
   ```bash
   # Update wrangler.toml with your account details
   account_id = "YOUR_ACCOUNT_ID"
   zone_id = "YOUR_ZONE_ID"
   
   # Update route pattern
   pattern = "yourdomain.com/*"
   ```

3. **Deploy**
   ```bash
   wrangler deploy
   ```

### Local Development

#### JavaScript (Recommended for production)
```bash
cd src/js
npm install wrangler -g
wrangler dev
```

#### Go
```bash
cd src/go
go run main.go ai_detection.go logging.go
# Server starts on http://localhost:8080
```

#### Python
```bash
cd src/python
python -m pip install aiohttp
python -m aiohttp.web -H localhost -P 8080 main:app
```

## 🤖 AI Crawler Detection

### Currently Detected Crawlers
- GPTBot (OpenAI)
- OAI-SearchBot (OpenAI)
- ChatGPT-User
- ClaudeBot (Anthropic)
- PerplexityBot
- Google-Extended
- Bytespider (ByteDance)
- Amazonbot
- Applebot-Extended
- CCBot (Common Crawl)

### AI-Optimized Content
When an AI crawler is detected, the worker serves:
```html
<html>
    <head><title>AI-Ready Content</title></head>
    <body>This is content optimized for AI crawlers.</body>
</html>
```

## 📊 Logging & Monitoring

### Request Logging
All requests are logged with:
- Timestamp and request details
- AI crawler detection status
- Response status codes
- Remote IP and referrer information

### Log Forwarding
Logs are forwarded to an external service (configurable in `logging.*` files):
```
LOG_SERVICE_URL = 'https://httpbin.org/post'
```

## 🔮 Future Functionality (In `extra/` folders)

### Advanced Metrics & Analytics
- Request volume tracking
- Response time monitoring  
- Geographic analysis
- User agent statistics
- Performance metrics

### Heuristic Bot Detection
Beyond user-agent matching:
- Suspicious header analysis
- Request rate monitoring
- Behavioral pattern detection
- Robots.txt compliance checking
- Asset request correlation

### Dynamic Content Processing
Strategy and Template patterns for:
- **Content Extraction**: Parse HTML/JSON/XML responses
- **Response Building**: Create summaries and structured data
- **Template Rendering**: Generate optimized output

## ⚙️ Configuration

### Environment Variables
- `ORIGIN_URL`: Target server for proxying (default: `http://www.mysite.com/`)
- `LOG_SERVICE_URL`: External logging endpoint

### Customization
- **AI Crawler List**: Update in `ai_detection.*` files
- **Content Optimization**: Modify AI-optimized content in `main.*` files
- **Logging Format**: Adjust log structure in `logging.*` files

## 🚀 Architecture Benefits

### Multi-Language Synchronization
- **Identical Functionality**: All implementations provide the same features
- **Language Choice Freedom**: Deploy with your preferred runtime
- **Easy Migration**: Switch between implementations without feature loss

### Modular Design
- **Separation of Concerns**: Each module has a single responsibility
- **Future-Proof**: Advanced features ready for activation
- **Maintainable**: Clean, well-documented codebase

### Performance Optimized
- **Lightweight**: Minimal overhead for regular traffic
- **Non-blocking**: Asynchronous logging doesn't impact response times
- **Efficient**: Smart caching and connection reuse

## 📚 Language-Specific Details

### JavaScript (Production Deployment)
- **Runtime**: Cloudflare Workers V8 engine
- **Performance**: Native execution, fastest cold starts
- **Deployment**: Direct deployment with `wrangler deploy`

### Go
- **Runtime**: Compiled to WebAssembly (WASM) via TinyGo
- **Performance**: Near-native performance with small binary size
- **Deployment**: Compile to WASM then deploy

### Python  
- **Runtime**: Pyodide (Python in WebAssembly)
- **Performance**: Full Python compatibility with moderate overhead
- **Deployment**: Experimental, use for development/testing

## 🤝 Contributing

### Adding New AI Crawlers
Update the crawler list in all three `ai_detection.*` files:
```javascript
// JavaScript
const AI_CRAWLERS = [..., "NewBot"];

// Go  
var aiCrawlers = []string{..., "NewBot"}

// Python
AI_CRAWLERS = [..., "NewBot"]
```

### Enabling Future Features
Move modules from `extra/` folders to main directories and update imports.

---

**Current Status**: ✅ Production ready with JavaScript implementation

**Future Roadmap**: 🔮 Advanced analytics and heuristic detection ready for activation
