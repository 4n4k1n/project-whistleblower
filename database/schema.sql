-- Users table for 42 students and staff
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    display_name TEXT NOT NULL,
    is_staff BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Reports table
CREATE TABLE IF NOT EXISTS reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reporter_id INTEGER NOT NULL,
    reported_student_login TEXT NOT NULL,
    project_name TEXT NOT NULL,
    reason TEXT NOT NULL,
    explanation TEXT NOT NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    reviewed_at DATETIME NULL,
    reviewed_by INTEGER NULL,
    FOREIGN KEY (reporter_id) REFERENCES users(id),
    FOREIGN KEY (reviewed_by) REFERENCES users(id)
);

-- Staff notifications for when reports reach threshold
CREATE TABLE IF NOT EXISTS staff_notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reported_student_login TEXT NOT NULL,
    project_name TEXT NOT NULL,
    report_count INTEGER NOT NULL,
    notification_sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    resolved BOOLEAN DEFAULT FALSE
);

-- Track false report patterns
CREATE TABLE IF NOT EXISTS user_report_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    total_reports INTEGER DEFAULT 0,
    approved_reports INTEGER DEFAULT 0,
    rejected_reports INTEGER DEFAULT 0,
    false_report_ratio REAL DEFAULT 0.0,
    warned BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Report reasons (predefined options)
CREATE TABLE IF NOT EXISTS report_reasons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reason TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL
);

-- Insert default report reasons
INSERT OR IGNORE INTO report_reasons (reason, description) VALUES 
('plagiarism', 'Code copied from another source without attribution'),
('collusion', 'Unauthorized collaboration between students'),
('external_help', 'Received unauthorized external assistance'),
('code_sharing', 'Sharing code with other students'),
('academic_dishonesty', 'Other forms of academic misconduct'),
('suspicious_similarity', 'Unusual similarities with other submissions');