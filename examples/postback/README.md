# Kite Connect Postback Server Example

This example demonstrates how to create a web server to receive and log postback notifications from Zerodha's Kite Connect API.

## What are Postbacks?

Postbacks (also known as webhooks) are HTTP POST requests that Kite Connect sends to your configured URL whenever there are order updates. This allows your application to receive real-time notifications about order status changes without having to continuously poll the API.

### When are postbacks sent?

Postbacks are sent for the following events:
- Order placement
- Order modification
- Order cancellation
- Order execution (partial or complete)
- Order rejection

## Quick Start

### 1. Build and Run the Server

```bash
# Navigate to the postback example directory
cd examples/postback

# Build the server
go build -o postback-server main.go

# Run the server
./postback-server
```

The server will start on port 8080 by default. You can customize the port using the `PORT` environment variable:

```bash
PORT=3000 ./postback-server
```

### 2. Configure Your Kite Connect App

1. Log in to your [Kite Connect Developer Console](https://developers.kite.trade/)
2. Navigate to your app settings
3. Set the postback URL to: `http://your-server-domain:8080/postback`
   - For local testing: `http://localhost:8080/postback`
   - For production: `https://yourdomain.com/postback`

### 3. Test the Server

Once running, you can:
- Visit `http://localhost:8080` in your browser to see server information
- Check health status at `http://localhost:8080/health`
- View logs in the console and in `postback_logs.json` file

## Server Endpoints

| Method | Endpoint   | Description |
|--------|------------|-------------|
| POST   | /postback  | Receives postback notifications from Kite Connect |
| GET    | /health    | Health check endpoint |
| GET    | /          | Server information page |

## Sample Postback Data

When Kite Connect sends a postback, it typically contains order information in JSON format. Here are some examples:

### Order Placed
```json
{
  "account_id": "XX1234",
  "order_id": "220224000000001",
  "exchange_order_id": "",
  "parent_order_id": null,
  "status": "OPEN",
  "status_message": "",
  "order_timestamp": "2022-02-24 09:15:01",
  "exchange_timestamp": "2022-02-24 09:15:01",
  "variety": "regular",
  "exchange": "NSE",
  "tradingsymbol": "RELIANCE",
  "instrument_token": 738561,
  "order_type": "LIMIT",
  "transaction_type": "BUY",
  "validity": "DAY",
  "product": "CNC",
  "quantity": 1,
  "disclosed_quantity": 0,
  "price": 2500.0,
  "trigger_price": 0,
  "average_price": 0,
  "filled_quantity": 0,
  "pending_quantity": 1,
  "cancelled_quantity": 0
}
```

### Order Executed
```json
{
  "account_id": "XX1234",
  "order_id": "220224000000001",
  "exchange_order_id": "1000000000000001",
  "parent_order_id": null,
  "status": "COMPLETE",
  "status_message": "",
  "order_timestamp": "2022-02-24 09:15:01",
  "exchange_timestamp": "2022-02-24 09:15:45",
  "variety": "regular",
  "exchange": "NSE",
  "tradingsymbol": "RELIANCE",
  "instrument_token": 738561,
  "order_type": "LIMIT",
  "transaction_type": "BUY",
  "validity": "DAY",
  "product": "CNC",
  "quantity": 1,
  "disclosed_quantity": 0,
  "price": 2500.0,
  "trigger_price": 0,
  "average_price": 2499.95,
  "filled_quantity": 1,
  "pending_quantity": 0,
  "cancelled_quantity": 0
}
```

## Configuration Options

### Environment Variables

- `PORT`: Server port (default: 8080)
- `LOG_FILE`: Log file path (default: postback_logs.json)

### Example with custom configuration:
```bash
PORT=3000 LOG_FILE=/var/log/kite-postbacks.json ./postback-server
```

## Building for Production

### Build for different platforms:

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o postback-server-linux main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o postback-server.exe main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o postback-server-macos main.go
```

### Docker Deployment

Create a `Dockerfile`:
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o postback-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/postback-server .
EXPOSE 8080
CMD ["./postback-server"]
```

Build and run:
```bash
docker build -t kite-postback-server .
docker run -p 8080:8080 kite-postback-server
```

## Security Considerations

### 1. Verify Postback Authenticity
In production, you should verify that postbacks are actually coming from Kite Connect. While Kite Connect doesn't currently provide signature verification, you can:
- Restrict access to your postback endpoint by IP (if Kite provides IP ranges)
- Use HTTPS for secure transmission
- Implement rate limiting

### 2. HTTPS in Production
Always use HTTPS in production to protect sensitive order data:
```bash
# Example with Let's Encrypt and reverse proxy
# Use nginx, Caddy, or similar to handle SSL termination
```

### 3. Input Validation
The server validates and parses incoming JSON data safely, but always review and test your error handling.

## Testing Postbacks Locally

### Using ngrok for local testing:
```bash
# Install ngrok (https://ngrok.com/)
# Start your postback server
./postback-server

# In another terminal, expose your local server
ngrok http 8080

# Use the ngrok HTTPS URL as your postback URL in Kite Connect app settings
# Example: https://abc123.ngrok.io/postback
```

### Manual Testing with curl:
```bash
# Test the postback endpoint
curl -X POST http://localhost:8080/postback \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "test123",
    "status": "COMPLETE",
    "tradingsymbol": "RELIANCE",
    "quantity": 1,
    "price": 2500.00
  }'

# Check health endpoint
curl http://localhost:8080/health
```

## Troubleshooting

### Common Issues:

1. **Port already in use**
   - Change the port: `PORT=3001 ./postback-server`

2. **Permission denied when creating log file**
   - Ensure write permissions: `chmod 755 .`
   - Or specify a different log file: `LOG_FILE=/tmp/postbacks.json ./postback-server`

3. **Postbacks not received**
   - Verify your postback URL is accessible from the internet
   - Check Kite Connect app configuration
   - Ensure the server is running and reachable
   - Check firewall settings

4. **SSL/TLS issues in production**
   - Use a reverse proxy like nginx or Caddy for SSL termination
   - Ensure your SSL certificate is valid

### Logs and Debugging:

The server logs all requests to both console and file. Check the logs for:
- Request headers and body
- Parsing errors
- Network issues

Example log entry:
```json
{
  "timestamp": "2022-02-24 09:15:45",
  "data": {
    "request_info": {
      "method": "POST",
      "url": "/postback",
      "headers": {
        "Content-Type": ["application/json"],
        "User-Agent": ["Kite-Connect/1.0"]
      },
      "content_type": "application/json",
      "user_agent": "Kite-Connect/1.0",
      "remote_addr": "127.0.0.1:54321",
      "body_size": 456
    },
    "postback": {
      "parsed_as": "kite_order",
      "order_data": {
        "order_id": "220224000000001",
        "status": "COMPLETE",
        ...
      },
      "raw_body": "{\"order_id\":\"220224000000001\",...}"
    }
  }
}
```

## Integration Examples

### Basic Order Processing:
```go
func processOrderUpdate(order kiteconnect.Order) {
    switch order.Status {
    case "COMPLETE":
        fmt.Printf("Order %s completed for %s\n", order.OrderID, order.TradingSymbol)
        // Update your database, send notifications, etc.
    case "REJECTED":
        fmt.Printf("Order %s rejected: %s\n", order.OrderID, order.StatusMessage)
        // Handle rejection, alert user, etc.
    }
}
```

### Database Integration:
```go
// Example with database storage
func saveOrderUpdate(db *sql.DB, order kiteconnect.Order) error {
    query := `
        INSERT INTO order_updates (order_id, status, timestamp, data) 
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE status=?, timestamp=?, data=?
    `
    orderJSON, _ := json.Marshal(order)
    _, err := db.Exec(query, 
        order.OrderID, order.Status, time.Now(), orderJSON,
        order.Status, time.Now(), orderJSON)
    return err
}
```

## Resources

- [Kite Connect API Documentation](https://kite.trade/docs/connect/v3/)
- [Kite Connect Developer Console](https://developers.kite.trade/)
- [GoKiteConnect GitHub Repository](https://github.com/zerodha/gokiteconnect)

## Support

For issues related to:
- **This example**: Create an issue in the GoKiteConnect repository
- **Kite Connect API**: Contact Zerodha support or check their documentation
- **Postback configuration**: Refer to Kite Connect documentation