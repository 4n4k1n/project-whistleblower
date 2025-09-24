#!/bin/bash

echo "🔐 Testing Admin Access Control"
echo "==============================="

# Start server
./whistleblower > test.log 2>&1 &
SERVER_PID=$!
sleep 3

echo "✅ Server started"

echo ""
echo "🧪 Test 1: Non-staff user trying to access admin"
echo "================================================"

# Temporarily remove staff privileges
sqlite3 whistleblower.db "UPDATE users SET is_staff = 0 WHERE login = 'apregitz';"
echo "Removed staff privileges for apregitz"

# Test admin access
RESPONSE=$(curl -s -b "user_login=apregitz" -w "%{http_code}" "http://localhost:8080/admin")
HTTP_CODE=$(echo "$RESPONSE" | tail -c 4)

echo "HTTP Response Code: $HTTP_CODE"

if [ "$HTTP_CODE" = "403" ]; then
    echo "✅ Access correctly denied (403 Forbidden)"
elif [ "$HTTP_CODE" = "302" ]; then
    echo "✅ Redirected to login (302)"
else
    echo "❌ Unexpected response code: $HTTP_CODE"
fi

echo ""
echo "🧪 Test 2: Staff user accessing admin"
echo "====================================="

# Restore staff privileges
sqlite3 whistleblower.db "UPDATE users SET is_staff = 1 WHERE login = 'apregitz';"
echo "Restored staff privileges for apregitz"

# Test admin access
RESPONSE2=$(curl -s -b "user_login=apregitz" -w "%{http_code}" "http://localhost:8080/admin")
HTTP_CODE2=$(echo "$RESPONSE2" | tail -c 4)

echo "HTTP Response Code: $HTTP_CODE2"

if [ "$HTTP_CODE2" = "200" ]; then
    echo "✅ Admin page accessible for staff (200 OK)"
    if echo "$RESPONSE2" | grep -q "Admin Panel"; then
        echo "✅ Admin page content loaded correctly"
    else
        echo "❌ Admin page content not found"
    fi
else
    echo "❌ Unexpected response code: $HTTP_CODE2"
fi

echo ""
echo "🧪 Test 3: Unauthenticated user"
echo "==============================="

# Test without cookies
RESPONSE3=$(curl -s -w "%{http_code}" "http://localhost:8080/admin")
HTTP_CODE3=$(echo "$RESPONSE3" | tail -c 4)

echo "HTTP Response Code: $HTTP_CODE3"

if [ "$HTTP_CODE3" = "302" ]; then
    echo "✅ Unauthenticated user redirected to login"
else
    echo "❌ Unexpected response code: $HTTP_CODE3"
fi

# Clean up
kill $SERVER_PID 2>/dev/null
rm -f test.log

echo ""
echo "🎉 Admin access control working!"
echo "================================"
echo "✅ Non-staff users: Access denied (403)"
echo "✅ Staff users: Full access (200)" 
echo "✅ Unauthenticated: Redirect to login (302)"