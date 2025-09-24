#!/bin/bash

echo "🗄️  Database Explorer for Whistleblower"
echo "======================================"

DB_FILE="whistleblower.db"

if [ ! -f "$DB_FILE" ]; then
    echo "❌ Database file not found: $DB_FILE"
    exit 1
fi

echo ""
echo "📊 Tables Overview:"
echo "==================="
sqlite3 $DB_FILE "SELECT name FROM sqlite_master WHERE type='table';"

echo ""
echo "👥 Users Table:"
echo "==============="
echo "Total users: $(sqlite3 $DB_FILE "SELECT COUNT(*) FROM users;")"
echo ""
echo "Sample users:"
sqlite3 $DB_FILE -header -column "SELECT login, display_name, email, is_staff FROM users LIMIT 10;"

echo ""
echo "📝 Reports Table:" 
echo "================="
echo "Total reports: $(sqlite3 $DB_FILE "SELECT COUNT(*) FROM reports;")"
if [ "$(sqlite3 $DB_FILE "SELECT COUNT(*) FROM reports;")" -gt 0 ]; then
    echo ""
    echo "Sample reports:"
    sqlite3 $DB_FILE -header -column "SELECT reporter_id, reported_student_login, project_name, reason, status FROM reports LIMIT 5;"
fi

echo ""
echo "🏷️  Report Reasons:"
echo "=================="
sqlite3 $DB_FILE -header -column "SELECT reason, description FROM report_reasons;"

echo ""
echo "🔧 Database Commands:"
echo "===================="
echo "Interactive shell:    sqlite3 $DB_FILE"
echo "View all users:       sqlite3 $DB_FILE 'SELECT * FROM users;'"
echo "Search users:         sqlite3 $DB_FILE \"SELECT * FROM users WHERE login LIKE '%search%';\""
echo "Count by campus:      sqlite3 $DB_FILE \"SELECT COUNT(*) FROM users WHERE email LIKE '%berlin%';\""
echo ""

echo "💡 Useful Queries:"
echo "=================="
echo "# Find German users"
echo "sqlite3 $DB_FILE \"SELECT login, display_name FROM users WHERE email LIKE '%berlin%' OR email LIKE '%wolfsburg%' LIMIT 10;\""
echo ""
echo "# Search for specific user"
echo "sqlite3 $DB_FILE \"SELECT * FROM users WHERE login = 'username';\""
echo ""
echo "# Get users with staff privileges"  
echo "sqlite3 $DB_FILE \"SELECT login, display_name FROM users WHERE is_staff = 1;\""