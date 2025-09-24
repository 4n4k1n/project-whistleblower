#!/bin/bash

echo "🏫 Finding 42 Heilbronn Campus ID"
echo "================================"

# Check if we have valid credentials
if [ -z "$OAUTH_42_CLIENT_ID" ] || [ -z "$OAUTH_42_CLIENT_SECRET" ]; then
    echo "⚠️  Please set your OAuth credentials first:"
    echo "export OAUTH_42_CLIENT_ID=your_client_id"
    echo "export OAUTH_42_CLIENT_SECRET=your_client_secret"
    echo ""
    echo "Or run: set -a; source .env; set +a"
    exit 1
fi

echo "🔑 Getting access token..."

# Get access token using client credentials
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

echo "✅ Access token obtained"
echo ""
echo "🔍 Searching for Heilbronn campus..."

# Get all campuses and search for Heilbronn
CAMPUS_RESPONSE=$(curl -s "https://api.intra.42.fr/v2/campus" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

echo "$CAMPUS_RESPONSE" | jq -r '.[] | select(.name | test("heilbronn|Heilbronn"; "i")) | "🏫 Campus: \(.name)\n   ID: \(.id)\n   Country: \(.country)\n   City: \(.city)\n"'

echo ""
echo "📋 All German campuses:"
echo "$CAMPUS_RESPONSE" | jq -r '.[] | select(.country == "Germany") | "   \(.name) (ID: \(.id)) - \(.city)"'

echo ""
echo "💡 If you don't see Heilbronn, try searching all campuses:"
echo "   curl -H \"Authorization: Bearer $ACCESS_TOKEN\" \"https://api.intra.42.fr/v2/campus\" | jq '.[] | {id, name, city, country}'"