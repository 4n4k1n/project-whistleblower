package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"whistleblower/models"
)

var (
	oauth2Config *oauth2.Config
)

func InitOAuth() {
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_42_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_42_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_42_REDIRECT_URL"),
		Scopes:       []string{"public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.intra.42.fr/oauth/authorize",
			TokenURL: "https://api.intra.42.fr/oauth/token",
		},
	}
}

func GetAuthURL(state string) string {
	return oauth2Config.AuthCodeURL(state)
}

func GetUserFromCode(code string) (*models.Auth42User, error) {
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://api.intra.42.fr/v2/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d", resp.StatusCode)
	}

	var user models.Auth42User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}

func SearchStudents(query string, token string) ([]models.StudentSearchResult, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.intra.42.fr/v2/users?search[login]=%s&per_page=10", query)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search students, status: %d", resp.StatusCode)
	}

	var users []struct {
		Login       string `json:"login"`
		DisplayName string `json:"displayname"`
		Email       string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	results := make([]models.StudentSearchResult, len(users))
	for i, user := range users {
		results[i] = models.StudentSearchResult{
			Login:       user.Login,
			DisplayName: user.DisplayName,
			Email:       user.Email,
		}
	}

	return results, nil
}

func GetStudentProjects(login string, token string) ([]string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.intra.42.fr/v2/users/%s/projects_users", login)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get student projects, status: %d", resp.StatusCode)
	}

	var projectUsers []struct {
		Project struct {
			Name string `json:"name"`
		} `json:"project"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&projectUsers); err != nil {
		return nil, err
	}

	projects := make([]string, len(projectUsers))
	for i, pu := range projectUsers {
		projects[i] = pu.Project.Name
	}

	return projects, nil
}

func GetCampusUsers(campusID int, token string, page int, perPage int) ([]models.Auth42User, error) {
	client := &http.Client{}
	// Use /v2/campus/{id}/users which works with client credentials
	url := fmt.Sprintf("https://api.intra.42.fr/v2/campus/%d/users?page=%d&per_page=%d", campusID, page, perPage)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get campus users, status: %d", resp.StatusCode)
	}

	var users []models.Auth42User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func GetAllCampusUsers(campusID int, token string) ([]models.Auth42User, error) {
	var allUsers []models.Auth42User
	page := 1
	perPage := 100

	for {
		users, err := GetCampusUsers(campusID, token, page, perPage)
		if err != nil {
			return nil, err
		}

		if len(users) == 0 {
			break
		}

		allUsers = append(allUsers, users...)

		if len(users) < perPage {
			break
		}

		page++
	}

	return allUsers, nil
}

func GetClientCredentialsToken() (string, error) {
	client := &http.Client{}
	
	data := fmt.Sprintf(`{
		"grant_type": "client_credentials",
		"client_id": "%s",
		"client_secret": "%s"
	}`, os.Getenv("OAUTH_42_CLIENT_ID"), os.Getenv("OAUTH_42_CLIENT_SECRET"))
	
	req, err := http.NewRequest("POST", "https://api.intra.42.fr/oauth/token", 
		strings.NewReader(data))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token, status: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}