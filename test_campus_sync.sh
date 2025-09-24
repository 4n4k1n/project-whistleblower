#!/bin/bash

echo "🧪 Testing Campus User Sync with Real Data"
echo "=========================================="

# Start server
./whistleblower > server.log 2>&1 &
SERVER_PID=$!
sleep 3

echo "✅ Server started (PID: $SERVER_PID)"

# Create a test user cookie to simulate login
echo "🍪 Creating test cookies..."
TEST_COOKIES="-b user_login=testuser"

echo ""
echo "🏫 Testing Campus 1 (Paris) - Small batch..."
RESPONSE_1=$(curl -s $TEST_COOKIES -X POST "http://localhost:8080/api/sync-users?campus_id=1")
echo "Paris response: $RESPONSE_1"

echo ""
echo "🏫 Testing Campus 51 (Berlin)..."
RESPONSE_51=$(curl -s $TEST_COOKIES -X POST "http://localhost:8080/api/sync-users?campus_id=51")
echo "Berlin response: $RESPONSE_51"

echo ""
echo "🏫 Testing Campus 44 (Wolfsburg)..."
RESPONSE_44=$(curl -s $TEST_COOKIES -X POST "http://localhost:8080/api/sync-users?campus_id=44")
echo "Wolfsburg response: $RESPONSE_44"

echo ""
echo "📊 Checking database..."
USER_COUNT=$(sqlite3 whistleblower.db "SELECT COUNT(*) FROM users;")
echo "Total users in database: $USER_COUNT"

if [ "$USER_COUNT" -gt 0 ]; then
    echo "✅ Users successfully synced!"
    echo "Sample users:"
    sqlite3 whistleblower.db "SELECT login, display_name FROM users LIMIT 5;"
else
    echo "❌ No users were synced"
fi

# Clean up
kill $SERVER_PID 2>/dev/null
rm -f server.log

echo ""
echo "🎉 Campus sync test complete!"