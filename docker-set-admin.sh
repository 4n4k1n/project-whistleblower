#!/bin/bash

# Script to set admin privileges for users in the containerized whistleblower app

if [ $# -eq 0 ]; then
    echo "Usage: $0 <username> [1|0]"
    echo "  username: 42 login name"
    echo "  1: make admin, 0: remove admin (default: 1)"
    echo ""
    echo "Examples:"
    echo "  $0 apregitz        # Make apregitz admin"
    echo "  $0 apregitz 0      # Remove admin from apregitz"
    echo ""
    echo "Current users:"
    docker-compose exec whistleblower sqlite3 /app/data/whistleblower.db "SELECT login, display_name, is_staff FROM users;"
    exit 1
fi

USERNAME=$1
IS_ADMIN=${2:-1}

echo "Setting admin status for '$USERNAME' to $IS_ADMIN..."

# Update the user in the database
docker-compose exec whistleblower sqlite3 /app/data/whistleblower.db "UPDATE users SET is_staff = $IS_ADMIN WHERE login = '$USERNAME';"

# Check if the update was successful
RESULT=$(docker-compose exec whistleblower sqlite3 /app/data/whistleblower.db "SELECT login, display_name, is_staff FROM users WHERE login = '$USERNAME';")

if [ -n "$RESULT" ]; then
    echo "Success! User status updated:"
    echo "$RESULT"
    echo ""
    echo "To access admin panel:"
    echo "1. Login via 42 OAuth at: http://localhost"
    echo "2. Then access: http://localhost/admin"
else
    echo "Error: User '$USERNAME' not found in database"
    echo "User must login via 42 OAuth first to be created in the database"
fi