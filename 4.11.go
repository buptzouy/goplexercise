package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const GitHubAPI = "https://api.github.com"

type Issue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	State string `json:"state,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-cli [create|read|update|delete] [args...]")
		return
	}

	action := os.Args[1]

	switch action {
	case "create":
		createIssue()
	case "read":
		if len(os.Args) != 3 {
			fmt.Println("Usage: github-cli read [issue_number]")
			return
		}
		readIssue(os.Args[2])
	case "update":
		if len(os.Args) != 3 {
			fmt.Println("Usage: github-cli update [issue_number]")
			return
		}
		updateIssue(os.Args[2])
	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: github-cli delete [issue_number]")
			return
		}
		deleteIssue(os.Args[2])
	default:
		fmt.Println("Unknown action:", action)
	}
}

func createIssue() {
	title, body := getEditorInput("Create a new issue")

	issue := &Issue{Title: title, Body: body}
	jsonData, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", GitHubAPI+"/repos/<owner>/<repo>/issues", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to create issue: %s\n", body)
	}

	fmt.Println("Issue created successfully!")
}

func readIssue(issueNumber string) {
	req, err := http.NewRequest("GET", GitHubAPI+"/repos/<owner>/<repo>/issues/"+issueNumber, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to read issue: %s\n", body)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %s\n", issue.Title)
	fmt.Printf("Body:\n%s\n", issue.Body)
}

func updateIssue(issueNumber string) {
	title, body := getEditorInput("Update the issue")

	issue := &Issue{Title: title, Body: body}
	jsonData, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", GitHubAPI+"/repos/<owner>/<repo>/issues/"+issueNumber, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to update issue: %s\n", body)
	}

	fmt.Println("Issue updated successfully!")
}

func deleteIssue(issueNumber string) {
	req, err := http.NewRequest("DELETE", GitHubAPI+"/repos/<owner>/<repo>/issues/"+issueNumber, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to delete issue: %s\n", body)
	}

	fmt.Println("Issue deleted successfully!")
}

func getEditorInput(prompt string) (string, string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpFile, err := os.CreateTemp("", "issue.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("# " + prompt + "\n\n")
	tmpFile.Close()

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	title := strings.TrimSpace(lines[0])
	body := strings.Join(lines[1:], "\n")

	return title, body
}
