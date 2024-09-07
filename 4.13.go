package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const OMDBAPIURL = "http://www.omdbapi.com/"

// Movie represents the JSON response from OMDb API
type Movie struct {
	Title  string `json:"Title"`
	Poster string `json:"Poster"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: poster <movie_name>")
		return
	}

	movieName := strings.Join(os.Args[1:], " ")
	apiKey := os.Getenv("OMDB_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the OMDB_API_KEY environment variable")
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("%s?apikey=%s&t=%s", OMDBAPIURL, apiKey, movieName)

	// Fetch movie details from OMDb API
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatalf("Failed to fetch movie data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		log.Fatalf("Failed to decode JSON response: %v", err)
	}

	if movie.Poster == "" || movie.Poster == "N/A" {
		log.Fatalf("No poster found for movie: %s", movieName)
	}

	// Download the poster image
	posterResp, err := http.Get(movie.Poster)
	if err != nil {
		log.Fatalf("Failed to download poster: %v", err)
	}
	defer posterResp.Body.Close()

	if posterResp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d when downloading poster", posterResp.StatusCode)
	}

	// Create a file to save the poster
	fileName := strings.ReplaceAll(movie.Title, " ", "_") + "_poster.jpg"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Save the poster image to the file
	_, err = io.Copy(file, posterResp.Body)
	if err != nil {
		log.Fatalf("Failed to save poster image: %v", err)
	}

	fmt.Printf("Poster downloaded successfully: %s\n", fileName)
}
