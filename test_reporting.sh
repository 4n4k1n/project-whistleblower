#!/bin/bash

echo "🧪 Testing Report Submission"
echo "============================"

# Start server
./whistleblower > test.log 2>&1 &
SERVER_PID=$!
sleep 3

echo "✅ Server started"

# Create test report data
REPORT_DATA='{
    "reported_student_login": "testuser", 
    "project_name": "libft",
    "reason": "plagiarism",
    "explanation": "Test report explanation - this is just a test"
}'

echo ""
echo "📝 Testing report submission with cookies..."
RESPONSE=$(curl -s -b "user_login=testuser" -X POST "http://localhost:8080/api/reports" \
    -H "Content-Type: application/json" \
    -d "$REPORT_DATA")

echo "Response: $RESPONSE"

if echo "$RESPONSE" | grep -q "Report submitted successfully\|report_id"; then
    echo "✅ Report submission working!"
    
    echo ""
    echo "📊 Checking database..."
    REPORT_COUNT=$(sqlite3 whistleblower.db "SELECT COUNT(*) FROM reports;")
    echo "Reports in database: $REPORT_COUNT"
    
    if [ "$REPORT_COUNT" -gt 0 ]; then
        echo "✅ Report saved to database!"
        echo ""
        echo "Recent report:"
        sqlite3 whistleblower.db -header -column "SELECT reported_student_login, project_name, reason, status FROM reports ORDER BY created_at DESC LIMIT 1;"
    fi
    
elif echo "$RESPONSE" | grep -q "Not authenticated\|User not found"; then
    echo "❌ Authentication issue - need proper OAuth login"
else
    echo "❌ Other issue with report submission"
fi

# Clean up
kill $SERVER_PID 2>/dev/null
rm -f test.log

echo ""
echo "💡 Summary:"
echo "==========="
echo "✅ Fixed: Project dropdown now always shows (with common projects)"
echo "✅ Fixed: Form validation should work properly now"
echo "🔧 Need: Proper OAuth login for full testing"
echo ""
echo "🚀 To test in browser:"
echo "1. Start: make dev"
echo "2. Login: http://localhost:8080/login"
echo "3. Try reporting: Select student → now project dropdown appears!"