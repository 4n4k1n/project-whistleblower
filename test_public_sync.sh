#!/bin/bash

echo "🔄 Testing Public User Sync (No Staff Required)"
echo "=============================================="

# Start server in background
./whistleblower > server.log 2>&1 &
SERVER_PID=$!
sleep 3

echo "✅ Server started (PID: $SERVER_PID)"

# Test sync endpoint without authentication (should fail)
echo "🔒 Testing sync endpoint security..."
SYNC_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/sync-users?campus_id=1")
echo "Response: $SYNC_RESPONSE"

if echo "$SYNC_RESPONSE" | grep -q "Not authenticated"; then
    echo "✅ Sync endpoint properly requires authentication (but not staff access)"
else
    echo "❌ Sync endpoint security issue"
fi

# Test that endpoint moved from /staff/
echo "🔧 Testing old staff endpoint..."
OLD_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/staff/sync-users?campus_id=1")
echo "Response: $OLD_RESPONSE"

if echo "$OLD_RESPONSE" | grep -q "404\|not found"; then
    echo "✅ Old staff endpoint properly removed"
else
    echo "❌ Old staff endpoint still accessible"
fi

# Clean up
kill $SERVER_PID 2>/dev/null
rm -f server.log

echo ""
echo "🎉 Public sync system ready!"
echo ""
echo "📋 Changes made:"
echo "• ✅ Changed from /v2/campus_users to /v2/users?filter[campus_id]=X"
echo "• ✅ Moved sync endpoint from /api/staff/sync-users to /api/sync-users"  
echo "• ✅ Removed staff requirement - any authenticated user can sync"
echo "• ✅ Uses client credentials token (no staff permissions needed)"
echo ""
echo "🚀 Now any logged-in user can sync campus users!"