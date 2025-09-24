package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"whistleblower/auth"
	"whistleblower/database"
	"whistleblower/models"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Login(c *gin.Context) {
	state := generateState()
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
	
	authURL := auth.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *Handler) Callback(c *gin.Context) {
	state := c.Query("state")
	storedState, err := c.Cookie("oauth_state")
	if err != nil || state != storedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	auth42User, err := auth.GetUserFromCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		return
	}

	user := &models.User{
		Login:       auth42User.Login,
		Email:       auth42User.Email,
		DisplayName: auth42User.DisplayName,
		IsStaff:     false,
	}

	if err := h.db.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token := generateToken()
	c.SetCookie("auth_token", token, 3600*24, "/", "", false, true)
	c.SetCookie("user_login", user.Login, 3600*24, "/", "", false, false)
	
	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}

func (h *Handler) SearchStudents(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter required"})
		return
	}

	_, err := c.Cookie("user_login")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Search in local database instead of 42 API for better performance
	results, err := h.db.SearchUsers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": results})
}

func (h *Handler) GetStudentProjects(c *gin.Context) {
	login := c.Param("login")
	
	token, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	projects, err := auth.GetStudentProjects(login, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *Handler) CreateReport(c *gin.Context) {
	userLogin, err := c.Cookie("user_login")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.db.GetUserByLogin(userLogin)
	if err != nil {
		// Debug: provide more info about the issue
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in database", 
			"login_from_cookie": userLogin,
			"hint": "Try logging out and logging back in via /login",
		})
		return
	}

	var req models.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report := &models.Report{
		ReporterID:          user.ID,
		ReportedStudentLogin: req.ReportedStudentLogin,
		ProjectName:         req.ProjectName,
		Reason:              req.Reason,
		Explanation:         req.Explanation,
	}

	if err := h.db.CreateReport(report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report"})
		return
	}

	reportCount, err := h.db.GetReportCountForProject(req.ReportedStudentLogin, req.ProjectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check report count"})
		return
	}

	if reportCount >= 3 {
		notification := &models.StaffNotification{
			ReportedStudentLogin: req.ReportedStudentLogin,
			ProjectName:         req.ProjectName,
			ReportCount:         reportCount,
		}
		
		if err := h.db.CreateStaffNotification(notification); err != nil {
			fmt.Printf("Failed to create staff notification: %v\n", err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Report submitted successfully",
		"report_id": report.ID,
	})
}

func (h *Handler) GetReportReasons(c *gin.Context) {
	reasons, err := h.db.GetReportReasons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report reasons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reasons": reasons})
}

func (h *Handler) GetPendingReports(c *gin.Context) {
	userLogin, err := c.Cookie("user_login")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.db.GetUserByLogin(userLogin)
	if err != nil || !user.IsStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "Staff access required"})
		return
	}

	reports, err := h.db.GetPendingReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func (h *Handler) ReviewReport(c *gin.Context) {
	userLogin, err := c.Cookie("user_login")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.db.GetUserByLogin(userLogin)
	if err != nil || !user.IsStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "Staff access required"})
		return
	}

	reportIDStr := c.Param("id")
	reportID, err := strconv.Atoi(reportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var req models.ReviewReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateReportStatus(reportID, req.Status, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report reviewed successfully"})
}

func (h *Handler) SyncCampusUsers(c *gin.Context) {
	// Simplify authentication - just check if user is logged in via cookie
	userLogin, err := c.Cookie("user_login")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated - please login first"})
		return
	}


	campusID := c.DefaultQuery("campus_id", "1")
	campusIDInt, err := strconv.Atoi(campusID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campus ID"})
		return
	}

	// Get OAuth token using client credentials (doesn't require user to be in DB)
	token, err := auth.GetClientCredentialsToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get OAuth token: " + err.Error()})
		return
	}

	auth42Users, err := auth.GetAllCampusUsers(campusIDInt, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campus users: " + err.Error()})
		return
	}

	users := make([]models.User, len(auth42Users))
	for i, auth42User := range auth42Users {
		users[i] = models.User{
			Login:       auth42User.Login,
			Email:       auth42User.Email,
			DisplayName: auth42User.DisplayName,
			IsStaff:     false,
		}
	}

	if err := h.db.BulkCreateUsers(users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save users to database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Users synced successfully from campus %d", campusIDInt),
		"count":   len(users),
		"requested_by": userLogin,
	})
}


func (h *Handler) GetUserStats(c *gin.Context) {
	count, err := h.db.GetUserCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_users": count,
	})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userLogin, err := c.Cookie("user_login")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.db.GetUserByLogin(userLogin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"login_from_cookie": userLogin,
			"in_database": false,
			"error": "User not in database - need to re-login via OAuth",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"login": user.Login,
		"display_name": user.DisplayName,
		"email": user.Email,
		"in_database": true,
	})
}

func (h *Handler) GetProjectStats(c *gin.Context) {
	userLogin, err := c.Cookie("user_login")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.db.GetUserByLogin(userLogin)
	if err != nil || !user.IsStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "Staff access required"})
		return
	}

	projectStats, err := h.db.GetMostReportedProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projectStats})
}

func (h *Handler) AdminPage(c *gin.Context) {
	userLogin, err := c.Cookie("user_login")
	if err != nil {
		// Redirect to login if not authenticated
		c.Redirect(302, "/login")
		return
	}

	user, err := h.db.GetUserByLogin(userLogin)
	if err != nil {
		// Redirect to login if user not found in database
		c.Redirect(302, "/login")
		return
	}

	if !user.IsStaff {
		// Show access denied page for non-staff users
		c.HTML(403, "access_denied.html", gin.H{
			"user_name": user.DisplayName,
			"message": "Access Denied: Staff privileges required to access the admin panel.",
		})
		return
	}

	// User is staff, show admin page
	c.HTML(200, "admin.html", gin.H{
		"user_name": user.DisplayName,
	})
}

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}