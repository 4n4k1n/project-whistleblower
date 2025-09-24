package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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