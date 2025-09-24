#!/bin/bash

echo "🔄 Testing User Sync Functionality"
echo "================================"

# Start server in background
./whistleblower > server.log 2>&1 &
SERVER_PID=$!
sleep 3

echo "✅ Server started (PID: $SERVER_PID)"

# Test user stats endpoint
echo "📊 Testing user stats..."
STATS_RESPONSE=$(curl -s http://localhost:8080/api/stats)
echo "Response: $STATS_RESPONSE"

if echo "$STATS_RESPONSE" | grep -q "total_users"; then
    echo "✅ User stats endpoint working"
else
    echo "❌ User stats endpoint failed"
fi

# Test admin page loads
echo "🔧 Testing admin panel..."
if curl -s http://localhost:8080/admin | grep -q "Admin Panel"; then
    echo "✅ Admin panel loads"
else
    echo "❌ Admin panel failed"
fi

# Test sync endpoint (will fail without proper auth, but should return error message)
echo "🔗 Testing sync endpoint..."
SYNC_RESPONSE=$(curl -s -X POST http://localhost:8080/api/staff/sync-users?campus_id=1)
echo "Response: $SYNC_RESPONSE"

if echo "$SYNC_RESPONSE" | grep -q "Not authenticated\|Staff access required"; then
    echo "✅ Sync endpoint properly protected"
else
    echo "❌ Sync endpoint security issue"
fi

# Check database
echo "🗄️ Checking database..."
USER_COUNT=$(sqlite3 whistleblower.db "SELECT COUNT(*) FROM users;")
echo "Current users in database: $USER_COUNT"

# Clean up
kill $SERVER_PID 2>/dev/null
rm -f server.log

echo ""
echo "🎉 User sync system ready!"
echo ""
echo "📋 How to use:"
echo "1. Login as staff user (set is_staff=1 in database)"
echo "2. Visit /admin to access the admin panel" 
echo "3. Click 'Sync Users' with desired campus ID"
echo "4. Users will be fetched from 42 API and saved to database"
echo ""
echo "🏫 Common Campus IDs:"
echo "- 1: Paris"
echo "- 7: Lyon"
echo "- 9: Brussels"
echo "- 12: Seoul"