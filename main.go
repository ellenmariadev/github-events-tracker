// Get the Pull Request Review Events for a user in the last 24 hours
//   - Load the github token from the .env file
//   - Load the github username from the .env file
//   - List the Pull Request Review Events for the user in the last 24 hours
//   - Print the Branch Name, Date and URL of the Pull Request Review Event

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if token == "" {
		fmt.Println("Missing github token environment variable")
		return
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		fmt.Println("Missing github token environment variable")
		return
	}

	// Load the Brazil/Sao Paulo time zone
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Get the current time
	now := time.Now().In(loc)

	// List events for the user
	events, _, err := client.Activity.ListEventsPerformedByUser(ctx, username, false, nil)
	if err != nil {
		fmt.Println("Error listing events:", err)
		return
	}

	for _, event := range events {
		localCreatedAt := event.CreatedAt.In(loc)
		// Check if the event was created in the last 24 hours
		if localCreatedAt.After(now.Add(-24 * time.Hour)) {
			if *event.Type == "PullRequestReviewEvent" {
				// Parse the Payload into a PullRequestReviewEvent
				var prReviewEvent github.PullRequestReviewEvent
				if err := json.Unmarshal(*event.RawPayload, &prReviewEvent); err != nil {
					fmt.Printf("Error parsing payload: %v\n", err)
					continue
				}
				// Print the PR URL
				fmt.Printf("Branch Name: %s | %s | URL: %s\n", *prReviewEvent.PullRequest.Head.Ref, localCreatedAt.Format("Mon, 02 Jan 2006 15:04:05"), *prReviewEvent.PullRequest.HTMLURL)
			}
		}
	}
}
