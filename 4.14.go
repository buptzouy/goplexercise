package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	GitHubAPIBaseURL = "https://api.github.com/repos/"
)

type Issue2 struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`
	State string `json:"state"`
}

type Milestone struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
}

type User struct {
	Login string `json:"login"`
}

var (
	issues     []Issue2
	milestones []Milestone
	users      []User
	once       sync.Once
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/issues", issuesHandler)
	http.HandleFunc("/milestones", milestonesHandler)
	http.HandleFunc("/users", usersHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchGitHubData(owner, repo, token string) {
	once.Do(func() {
		// Fetch Issues
		issuesURL := fmt.Sprintf("%s%s/%s/issues", GitHubAPIBaseURL, owner, repo)
		issues = fetchData[Issue2](issuesURL, token)

		// Fetch Milestones
		milestonesURL := fmt.Sprintf("%s%s/%s/milestones", GitHubAPIBaseURL, owner, repo)
		milestones = fetchData[Milestone](milestonesURL, token)

		// Fetch Users (Contributors)
		usersURL := fmt.Sprintf("%s%s/%s/contributors", GitHubAPIBaseURL, owner, repo)
		users = fetchData[User](usersURL, token)
	})
}

func fetchData[T any](url, token string) []T {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch data from GitHub: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Error: received status code %d, body: %s", resp.StatusCode, body)
	}

	var data []T
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalf("Failed to decode JSON response: %v", err)
	}
	return data
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
		<html>
		<head><title>GitHub Navigation</title></head>
		<body>
			<h1>GitHub Data Navigation</h1>
			<ul>
				<li><a href="/issues">Issues</a></li>
				<li><a href="/milestones">Milestones</a></li>
				<li><a href="/users">Users</a></li>
			</ul>
		</body>
		</html>`
	fmt.Fprint(w, tmpl)
}

func issuesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<html>
		<head><title>Issues</title></head>
		<body>
			<h1>Issues</h1>
			<ul>
				{{range .}}
					<li>#{{.Number}}: {{.Title}} ({{.User.Login}}) - {{.State}}</li>
				{{end}}
			</ul>
			<a href="/">Back</a>
		</body>
		</html>`
	renderTemplate(w, tmpl, issues)
}

func milestonesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<html>
		<head><title>Milestones</title></head>
		<body>
			<h1>Milestones</h1>
			<ul>
				{{range .}}
					<li>{{.Title}} - {{.Description}} ({{.State}})</li>
				{{end}}
			</ul>
			<a href="/">Back</a>
		</body>
		</html>`
	renderTemplate(w, tmpl, milestones)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<html>
		<head><title>Users</title></head>
		<body>
			<h1>Users</h1>
			<ul>
				{{range .}}
					<li>{{.Login}}</li>
				{{end}}
			</ul>
			<a href="/">Back</a>
		</body>
		</html>`
	renderTemplate(w, tmpl, users)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t := template.Must(template.New("webpage").Parse(tmpl))
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
