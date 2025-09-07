#!/bin/bash

# Test script for Kite Connect Postback Server
echo "🧪 Testing Kite Connect Postback Server"
echo "========================================"

# Check if server is running
SERVER_URL="http://localhost:8080"
echo "1. Checking if server is running at $SERVER_URL..."

if curl -s "$SERVER_URL/health" > /dev/null; then
    echo "✅ Server is running"
else
    echo "❌ Server is not running. Please start the server first with: go run main.go"
    exit 1
fi

echo ""
echo "2. Testing health endpoint..."
curl -s "$SERVER_URL/health" | jq '.' 2>/dev/null || curl -s "$SERVER_URL/health"

echo ""
echo ""
echo "3. Testing postback endpoint with sample order data..."

# Sample order placed
echo "📤 Sending sample 'Order Placed' postback..."
curl -X POST "$SERVER_URL/postback" \
  -H "Content-Type: application/json" \
  -H "User-Agent: Kite-Connect/1.0" \
  -d '{
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
  }' 2>/dev/null

echo ""
echo ""

# Sample order executed
echo "📤 Sending sample 'Order Executed' postback..."
curl -X POST "$SERVER_URL/postback" \
  -H "Content-Type: application/json" \
  -H "User-Agent: Kite-Connect/1.0" \
  -d '{
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
  }' 2>/dev/null

echo ""
echo ""

# Test with invalid JSON
echo "📤 Sending invalid JSON to test error handling..."
curl -X POST "$SERVER_URL/postback" \
  -H "Content-Type: application/json" \
  -d '{"invalid": json}' 2>/dev/null

echo ""
echo ""

# Test with empty body
echo "📤 Sending empty body..."
curl -X POST "$SERVER_URL/postback" \
  -H "Content-Type: application/json" \
  -d '' 2>/dev/null

echo ""
echo ""

# Test wrong method
echo "📤 Testing wrong HTTP method (should fail)..."
curl -X GET "$SERVER_URL/postback" 2>/dev/null

echo ""
echo ""
echo "✅ Tests completed!"
echo ""
echo "📋 Check the server logs in the console and in postback_logs.json"
echo "🌐 Visit $SERVER_URL in your browser for server information"