#!/bin/bash

sqlite3 whistleblower.db "UPDATE users SET is_staff = 1 WHERE login ='apregitz';"