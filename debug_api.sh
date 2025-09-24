#!/bin/bash

echo "🔍 Debug Campus Users API Call"
echo "=============================="

# Load environment
set -a; source .env; set +a

echo "🔑 Getting OAuth token..."
TOKEN_RESPONSE=$(curl -s -X POST "https://api.intra.42.fr/oauth/token" \
    -H "Content-Type: application/json" \
    -d "{
        \"grant_type\": \"client_credentials\",
        \"client_id\": \"$OAUTH_42_CLIENT_ID\",
        \"client_secret\": \"$OAUTH_42_CLIENT_SECRET\"
    }")

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "❌ Failed to get access token"
    echo "Response: $TOKEN_RESPONSE"
    exit 1
fi

echo "✅ Token obtained: ${ACCESS_TOKEN:0:20}..."

# Test different campus endpoints
echo ""
echo "🏫 Testing Campus 51 (Berlin)..."
RESPONSE_51=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" \
    "https://api.intra.42.fr/v2/users?filter[campus_id]=51&per_page=5")

echo "Response length: ${#RESPONSE_51} characters"
echo "First few users:"
echo "$RESPONSE_51" | head -c 500
echo ""

echo "🏫 Testing Campus 44 (Wolfsburg)..."
RESPONSE_44=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" \
    "https://api.intra.42.fr/v2/users?filter[campus_id]=44&per_page=5")

echo "Response length: ${#RESPONSE_44} characters"
echo "First few users:"
echo "$RESPONSE_44" | head -c 500
echo ""

echo "🏫 Testing Campus 1 (Paris)..."
RESPONSE_1=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" \
    "https://api.intra.42.fr/v2/users?filter[campus_id]=1&per_page=5")

echo "Response length: ${#RESPONSE_1} characters"
echo "First few users:"
echo "$RESPONSE_1" | head -c 500
echo ""

echo "💡 If responses are empty, the campus filter might not work as expected"