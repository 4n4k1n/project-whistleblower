#!/bin/bash

echo "🧪 Testing Whistleblower Application"
echo "=================================="

# Set test environment variables
export OAUTH_42_CLIENT_ID=test_client
export OAUTH_42_CLIENT_SECRET=test_secret  
export OAUTH_42_REDIRECT_URL=http://localhost:8080/callback
export PORT=8080

echo "✅ Environment variables set"

# Test build
echo "🔨 Testing build..."
go build -o whistleblower . 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

# Test database initialization
echo "🗄️  Testing database..."
if [ -f "whistleblower.db" ]; then
    echo "✅ Database exists"
    
    # Check report reasons
    REASON_COUNT=$(sqlite3 whistleblower.db "SELECT COUNT(*) FROM report_reasons;")
    if [ "$REASON_COUNT" -eq "6" ]; then
        echo "✅ Report reasons loaded ($REASON_COUNT)"
    else
        echo "❌ Wrong number of report reasons: $REASON_COUNT"
    fi
    
    # Check tables exist
    TABLES=$(sqlite3 whistleblower.db ".tables" | wc -w)
    if [ "$TABLES" -eq "5" ]; then
        echo "✅ All database tables exist ($TABLES)"
    else
        echo "❌ Missing database tables: $TABLES/5"
    fi
else
    echo "❌ Database not found"
    exit 1
fi

# Start server in background and test endpoints
echo "🌐 Testing server..."
timeout 3s ./whistleblower > server.log 2>&1 &
SERVER_PID=$!
sleep 2

# Test homepage
if curl -s http://localhost:8080/ | grep -q "42 Academic Integrity Portal"; then
    echo "✅ Homepage loads"
else
    echo "❌ Homepage failed"
fi

# Test API endpoint
if curl -s http://localhost:8080/api/report-reasons | grep -q "plagiarism"; then
    echo "✅ API endpoint works"
else
    echo "❌ API endpoint failed"
fi

# Test OAuth redirect
if curl -s http://localhost:8080/login | grep -q "api.intra.42.fr"; then
    echo "✅ OAuth redirect works"
else
    echo "❌ OAuth redirect failed"
fi

# Clean up
kill $SERVER_PID 2>/dev/null || true
rm -f server.log

echo ""
echo "🎉 All tests completed!"
echo ""
echo "📋 Next Steps:"
echo "1. Set up your 42 OAuth app at: https://profile.intra.42.fr/oauth/applications"
echo "2. Copy .env.example to .env and add your credentials"
echo "3. Run: make dev"
echo "4. Visit: http://localhost:8080"