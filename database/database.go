package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"whistleblower/models"
)

type DB struct {
	*sql.DB
}

func NewDatabase(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	
	if err := database.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

func (db *DB) InitSchema() error {
	schemaPath := filepath.Join("database", "schema.sql")
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

func (db *DB) CreateUser(user *models.User) error {
	query := `INSERT OR REPLACE INTO users (login, email, display_name, is_staff) 
			  VALUES (?, ?, ?, ?)`
	
	result, err := db.Exec(query, user.Login, user.Email, user.DisplayName, user.IsStaff)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (db *DB) GetUserByLogin(login string) (*models.User, error) {
	query := `SELECT id, login, email, display_name, is_staff, created_at FROM users WHERE login = ?`
	
	var user models.User
	err := db.QueryRow(query, login).Scan(
		&user.ID, &user.Login, &user.Email, 
		&user.DisplayName, &user.IsStaff, &user.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (db *DB) CreateReport(report *models.Report) error {
	query := `INSERT INTO reports (reporter_id, reported_student_login, project_name, reason, explanation) 
			  VALUES (?, ?, ?, ?, ?)`
	
	result, err := db.Exec(query, report.ReporterID, report.ReportedStudentLogin, 
		report.ProjectName, report.Reason, report.Explanation)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	report.ID = int(id)
	return nil
}

func (db *DB) GetReportCountForProject(studentLogin, projectName string) (int, error) {
	query := `SELECT COUNT(*) FROM reports WHERE reported_student_login = ? AND project_name = ? AND status = 'pending'`
	
	var count int
	err := db.QueryRow(query, studentLogin, projectName).Scan(&count)
	return count, err
}

func (db *DB) GetPendingReports() ([]models.Report, error) {
	query := `SELECT id, reporter_id, reported_student_login, project_name, reason, explanation, status, created_at 
			  FROM reports WHERE status = 'pending' ORDER BY created_at DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var report models.Report
		err := rows.Scan(&report.ID, &report.ReporterID, &report.ReportedStudentLogin,
			&report.ProjectName, &report.Reason, &report.Explanation, 
			&report.Status, &report.CreatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (db *DB) UpdateReportStatus(reportID int, status string, reviewerID int) error {
	query := `UPDATE reports SET status = ?, reviewed_by = ?, reviewed_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, status, reviewerID, reportID)
	return err
}

func (db *DB) GetReportReasons() ([]models.ReportReason, error) {
	query := `SELECT id, reason, description FROM report_reasons ORDER BY reason`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reasons []models.ReportReason
	for rows.Next() {
		var reason models.ReportReason
		err := rows.Scan(&reason.ID, &reason.Reason, &reason.Description)
		if err != nil {
			return nil, err
		}
		reasons = append(reasons, reason)
	}

	return reasons, nil
}

func (db *DB) CreateStaffNotification(notification *models.StaffNotification) error {
	query := `INSERT INTO staff_notifications (reported_student_login, project_name, report_count) 
			  VALUES (?, ?, ?)`
	
	result, err := db.Exec(query, notification.ReportedStudentLogin, 
		notification.ProjectName, notification.ReportCount)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	notification.ID = int(id)
	return nil
}

func (db *DB) UpdateUserReportStats(userID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var totalReports, approvedReports, rejectedReports int
	
	query := `SELECT 
		COUNT(*) as total,
		COUNT(CASE WHEN status = 'approved' THEN 1 END) as approved,
		COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected
		FROM reports WHERE reporter_id = ? AND status != 'pending'`
	
	err = tx.QueryRow(query, userID).Scan(&totalReports, &approvedReports, &rejectedReports)
	if err != nil {
		return err
	}

	var falseReportRatio float64
	if totalReports > 0 {
		falseReportRatio = float64(rejectedReports) / float64(totalReports)
	}

	upsertQuery := `INSERT OR REPLACE INTO user_report_stats 
		(user_id, total_reports, approved_reports, rejected_reports, false_report_ratio)
		VALUES (?, ?, ?, ?, ?)`
	
	_, err = tx.Exec(upsertQuery, userID, totalReports, approvedReports, rejectedReports, falseReportRatio)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) BulkCreateUsers(users []models.User) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO users (login, email, display_name, is_staff) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, user := range users {
		_, err := stmt.Exec(user.Login, user.Email, user.DisplayName, user.IsStaff)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db *DB) GetUserCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (db *DB) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, login, email, display_name, is_staff, created_at FROM users ORDER BY login`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Login, &user.Email, &user.DisplayName, &user.IsStaff, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}