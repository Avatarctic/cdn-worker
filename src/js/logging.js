// Logging Module
// Contains log entry structure and log forwarding functionality

// Log service configuration
const LOG_SERVICE_URL = 'https://httpbin.org/post';

// Forward log data to the configured log service
async function forwardLog(logEntry) {
  try {
    const response = await fetch(LOG_SERVICE_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'Cloudflare-Worker-Logger/1.0',
      },
      body: JSON.stringify(logEntry),
    });

    if (!response.ok) {
      console.error(`Log service returned non-success status: ${response.status}`);
      const errorText = await response.text();
      console.error(`Log service error response: ${errorText}`);
    }
  } catch (error) {
    console.error('Failed to forward log:', error);
  }
}

// Export functions
export { forwardLog };
