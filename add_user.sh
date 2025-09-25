#!/bin/bash

sqlite3 whistleblower.db "INSERT INTO users (login, email, display_name, is_staff) VALUES ('new_user_login', 'new_user_email@example.com', 'New User Name', 0);"