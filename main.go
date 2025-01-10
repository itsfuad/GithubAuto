package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"githubauto/colors"
)

// GitHub API URL
const apiURL = "https://api.github.com"

// Store user token for authentication
var token string

// GitHub repository structure to unmarshal response
type Repository struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
	Language      string `json:"language"`
	DefaultBranch string `json:"default_branch"`
	License       struct {
		Name string `json:"name"`
	} `json:"license"`
	Size   int `json:"size"`
	Stars  int `json:"stargazers_count"`
	Forks  int `json:"forks_count"`
	Url	string `json:"html_url"`
}

// Constant for repeated literal
const repoNameFormat = "Repository Name: %s\n"
const descriptionFormat = "Description: %s\n"
const urlFormat = "URL: %s\n"


// GitHub Issue structure
type Issue struct {
	Title string `json:"title"`
	State string `json:"state"`
	URL   string `json:"html_url"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
	CreatedAt string `json:"created_at"`
}

// Notification structure for alerting
type Notification struct {
	ID      string `json:"id"`
	Subject struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"subject"`
	UpdatedAt string `json:"updated_at"`
}


func main() {
	// Define flags for commands and options
	saveTokenFlag := flag.Bool("save-token", false, "Save GitHub Token")
	searchRepoFlag := flag.String("search-repo", "", "Search for a GitHub repository")
	showRepoFlag := flag.String("show-repo", "", "Show details of a GitHub repository")
	queryIssuesFlag := flag.String("query-issues", "", "Query issues for a specific repository")
	notifyFlag := flag.Bool("notify", false, "Check notifications for the user")
	allRepoFlag := flag.Bool("all-repo", false, "Fetch all repositories for the authenticated user")


	// Set usage for the flags in color
	flag.Usage = func() {
		colors.YELLOW.Println("Usage: githubcli [options]")
		colors.YELLOW.Println("Options:")
		colors.CYAN.Println("  -save-token\t\tSave GitHub Token")
		colors.CYAN.Println("  -all-repo\t\tFetch all repositories for the authenticated user")
		colors.CYAN.Println("  -search-repo <repo>\t\tSearch for a GitHub repository")
		colors.CYAN.Println("  -show-repo <repo>\t\tShow details of a GitHub repository")
		colors.CYAN.Println("  -query-issues\t\tQuery issues for a specific repository")
		colors.CYAN.Println("  -notify\t\tCheck notifications for the user")
	}

	// Parse the flags
	flag.Parse()

	// Handle saving the GitHub token
	if *saveTokenFlag {
		saveGitHubToken()
		return
	}

	// Handle fetching all repositories
	if *allRepoFlag {
		fetchAllRepositories()
		return
	}	

	// Handle searching for a repository
	if *searchRepoFlag != "" {
		searchRepository(*searchRepoFlag)
		return
	}

	// Handle showing repository details
	if *showRepoFlag != "" {
		showRepository(*showRepoFlag)
		return
	}

	// Handle querying issues in a repo
	if *queryIssuesFlag != "" {
		queryIssues(*queryIssuesFlag)
		return
	}

	// Handle notifications
	if *notifyFlag {
		checkNotifications()
	}
	// If no flags are provided, show usage
	if flag.NFlag() == 0 {
		flag.Usage()
	}
}

// Save GitHub Token to file
func saveGitHubToken() {
	colors.YELLOW.Println("Enter your GitHub Personal Access Token: ")
	var inputToken string
	_, err := fmt.Scanln(&inputToken)
	if err != nil {
		log.Fatalf("Error reading token: %v\n", err)
	}
	// Save token to a file
	err = os.WriteFile("github_token.txt", []byte(inputToken), 0600)
	if err != nil {
		log.Fatalf("Failed to save token: %v\n", err)
	}
	colors.GREEN.Println("Token saved successfully!")
}

// Load GitHub Token from file
func loadGitHubToken() string {
	data, err := os.ReadFile("github_token.txt")
	if err != nil {
		//if file does not exist, prompt user to save token
		if os.IsNotExist(err) {
			colors.YELLOW.Print("GitHub Token not found. Please save your token using the ")
			colors.GREY.Print("-save-token")
			colors.YELLOW.Println(" flag.")
			os.Exit(1)
		}
	}
	return string(data)
}

// Make an authenticated request to the GitHub API
func makeGitHubRequest(url string) (*http.Response, error) {
	token = loadGitHubToken() // Load token from file
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v\n", err)
	}

	// Set Authorization header with token
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	return client.Do(req)
}

// Search for GitHub repository by name
func searchRepository(repoName string) {
	url := fmt.Sprintf("%s/search/repositories?q=%s", apiURL, repoName)
	resp, err := makeGitHubRequest(url)
	if err != nil {
		log.Fatalf("Error searching repository: %v\n", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing response: %v\n", err)
	}

	items := result["items"].([]interface{})
	for _, item := range items {
		repo := item.(map[string]interface{})
		colors.BLUE.Printf(repoNameFormat, repo["name"])
		colors.CYAN.Printf(descriptionFormat, repo["description"])
		colors.GREEN.Printf("Owner: %s\n", repo["owner"].(map[string]interface{})["login"])
		colors.MAGENTA.Printf("Stars: %v\n", repo["stargazers_count"])
		colors.YELLOW.Printf("Forks: %v\n\n", repo["forks_count"])
		colors.CYAN.Printf(urlFormat, repo["html_url"])
	}
}

// Show details of a specific repository
func showRepository(repoName string) {
	url := fmt.Sprintf("%s/repos/%s", apiURL, repoName)
	resp, err := makeGitHubRequest(url)
	if err != nil {
		log.Fatalf("Error showing repository: %v\n", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("GitHub API returned an error: %s\nDetails: %s", resp.Status, string(body))
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
	colors.BLUE.Printf(repoNameFormat, repo.Name)
	}

	colors.BLUE.Printf(repoNameFormat, repo.Name)
	colors.CYAN.Printf(descriptionFormat, repo.Description)
	colors.GREEN.Printf("Owner: %s\n", repo.Owner.Login)
	colors.YELLOW.Printf("Language: %s\n", repo.Language)
	colors.MAGENTA.Printf("Default Branch: %s\n", repo.DefaultBranch)
	colors.PURPLE.Printf("License: %s\n", repo.License.Name)
	colors.GREY.Printf("Size: %d KB\n", repo.Size)
	colors.RED.Printf("Stars: %d\n", repo.Stars)
	colors.CYAN.Printf("Forks: %d\n", repo.Forks)
	colors.GREEN.Printf(urlFormat, repo.Url)
}


// Query issues in a specific repository
func queryIssues(repoName string) {
	url := fmt.Sprintf("%s/repos/%s/issues", apiURL, repoName)
	resp, err := makeGitHubRequest(url)
	if err != nil {
		log.Fatalf("Error querying issues: %v\n", err)
	}
	defer resp.Body.Close()

	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		log.Fatalf("Error parsing issues: %v\n", err)
	}

	for _, issue := range issues {
		colors.YELLOW.Printf("Title: %s\n", issue.Title)
		colors.RED.Printf("State: %s\n", issue.State)
		colors.GREEN.Printf(urlFormat, issue.URL)
		colors.CYAN.Printf("Labels: %v\n", issue.Labels)
		colors.MAGENTA.Printf("Created At: %s\n\n", issue.CreatedAt)
	}
}

func fetchAllRepositories() {
	url := fmt.Sprintf("%s/user/repos", apiURL)
	resp, err := makeGitHubRequest(url)
	if err != nil {
		log.Fatalf("Error fetching repositories: %v\n", err)
	}
	defer resp.Body.Close()

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Fatalf("Error parsing repositories: %v\n", err)
	}

	languages := make(map[string]int)

	for _, repo := range repos {
		colors.BLUE.Printf(repoNameFormat, repo.Name)
		colors.CYAN.Printf(descriptionFormat, repo.Description)
		colors.GREEN.Printf("Language: %s\n", repo.Language)
		colors.MAGENTA.Printf("Stars: %d\n", repo.Stars)
		colors.YELLOW.Printf("Forks: %d\n", repo.Forks)
		colors.GREY.Printf("Size: %d KB\n", repo.Size)
		colors.CYAN.Printf("URL: %s\n\n", repo.Url)

		// Count the number of repositories per language
		languages[repo.Language]++
	}

	if len(repos) > 0 {
		if len(repos) == 1 {
			colors.GREEN.Println("1 repository found!")
		} else {
			colors.GREEN.Printf("%d repositories found!\n", len(repos))
		}

		// Display the number of repositories per language
		colors.YELLOW.Println("Languages:")
		for lang, count := range languages {
			colors.CYAN.Printf("  %s: %d\n", lang, count)
		}
	} else {
		colors.RED.Println("No repositories found!")
	}
}

// Check for GitHub notifications
func checkNotifications() {
	url := fmt.Sprintf("%s/notifications", apiURL)
	resp, err := makeGitHubRequest(url)
	if err != nil {
		log.Fatalf("Error checking notifications: %v\n", err)
	}
	defer resp.Body.Close()

	var notifications []Notification
	if err := json.NewDecoder(resp.Body).Decode(&notifications); err != nil {
		log.Fatalf("Error parsing notifications: %v\n", err)
	}

	for _, notification := range notifications {
		colors.PURPLE.Printf("Notification Title: %s\n", notification.Subject.Title)
		colors.CYAN.Printf(urlFormat, notification.Subject.URL)
		colors.GREY.Printf("Updated At: %s\n\n", notification.UpdatedAt)
	}
}