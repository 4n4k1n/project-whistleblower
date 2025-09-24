package models

import (
	"time"
)

type User struct {
	ID          int       `json:"id" db:"id"`
	Login       string    `json:"login" db:"login"`
	Email       string    `json:"email" db:"email"`
	DisplayName string    `json:"display_name" db:"display_name"`
	IsStaff     bool      `json:"is_staff" db:"is_staff"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Report struct {
	ID                   int        `json:"id" db:"id"`
	ReporterID          int        `json:"reporter_id" db:"reporter_id"`
	ReportedStudentLogin string     `json:"reported_student_login" db:"reported_student_login"`
	ProjectName         string     `json:"project_name" db:"project_name"`
	Reason              string     `json:"reason" db:"reason"`
	Explanation         string     `json:"explanation" db:"explanation"`
	Status              string     `json:"status" db:"status"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	ReviewedAt          *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewedBy          *int       `json:"reviewed_by,omitempty" db:"reviewed_by"`
}

type StaffNotification struct {
	ID                   int       `json:"id" db:"id"`
	ReportedStudentLogin string    `json:"reported_student_login" db:"reported_student_login"`
	ProjectName         string    `json:"project_name" db:"project_name"`
	ReportCount         int       `json:"report_count" db:"report_count"`
	NotificationSentAt  time.Time `json:"notification_sent_at" db:"notification_sent_at"`
	Resolved            bool      `json:"resolved" db:"resolved"`
}

type UserReportStats struct {
	ID                int     `json:"id" db:"id"`
	UserID            int     `json:"user_id" db:"user_id"`
	TotalReports      int     `json:"total_reports" db:"total_reports"`
	ApprovedReports   int     `json:"approved_reports" db:"approved_reports"`
	RejectedReports   int     `json:"rejected_reports" db:"rejected_reports"`
	FalseReportRatio  float64 `json:"false_report_ratio" db:"false_report_ratio"`
	Warned            bool    `json:"warned" db:"warned"`
}

type ReportReason struct {
	ID          int    `json:"id" db:"id"`
	Reason      string `json:"reason" db:"reason"`
	Description string `json:"description" db:"description"`
}

type CreateReportRequest struct {
	ReportedStudentLogin string `json:"reported_student_login" binding:"required"`
	ProjectName         string `json:"project_name" binding:"required"`
	Reason              string `json:"reason" binding:"required"`
	Explanation         string `json:"explanation" binding:"required"`
}

type ReviewReportRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

type StudentSearchResult struct {
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Projects    []string `json:"projects,omitempty"`
}

type Auth42User struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	DisplayName string `json:"displayname"`
}