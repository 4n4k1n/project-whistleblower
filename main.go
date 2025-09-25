package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"whistleblower/auth"
	"whistleblower/database"
	"whistleblower/handlers"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatal("Failed to load environment variables:", err)
	}

	auth.InitOAuth()

	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "whistleblower.db"
	}
	
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	h := handlers.NewHandler(db)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	r.GET("/login", h.Login)
	r.GET("/callback", h.Callback)
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", gin.H{})
	})

	r.GET("/admin", h.AdminPage)

	api := r.Group("/api")
	{
		api.GET("/students/search", h.SearchStudents)
		api.GET("/students/:login/projects", h.GetStudentProjects)
		api.POST("/reports", h.CreateReport)
		api.GET("/report-reasons", h.GetReportReasons)
		api.GET("/stats", h.GetUserStats)
		api.GET("/me", h.GetCurrentUser) // Debug endpoint
		api.POST("/sync-users", h.SyncCampusUsers) // Moved from staff-only
		
		staff := api.Group("/staff")
		{
			staff.GET("/reports", h.GetPendingReports)
			staff.PUT("/reports/:id", h.ReviewReport)
			staff.GET("/project-stats", h.GetProjectStats)
			staff.POST("/bulk-project-action", h.BulkProjectAction)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	requiredEnvVars := []string{
		"OAUTH_42_CLIENT_ID",
		"OAUTH_42_CLIENT_SECRET", 
		"OAUTH_42_REDIRECT_URL",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Printf("Warning: %s environment variable not set", envVar)
		}
	}

	return nil
}