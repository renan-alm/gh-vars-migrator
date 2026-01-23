package main

import (
	"fmt"
	"log"

	"github.com/cli/go-gh/v2/pkg/api"
)

func main() {
	fmt.Println("gh-vars-migrator - GitHub CLI extension for variables migration")
	
	// Create a GitHub API client using the authenticated user's token
	client, err := api.DefaultRESTClient()
	if err != nil {
		log.Fatalf("Failed to create GitHub API client: %v", err)
	}

	// Fetch the authenticated user info
	response := struct {
		Login string `json:"login"`
		Name  string `json:"name"`
	}{}

	err = client.Get("user", &response)
	if err != nil {
		log.Fatalf("Failed to fetch user info: %v", err)
	}

	fmt.Printf("\nAuthenticated as: %s (%s)\n", response.Login, response.Name)
}
