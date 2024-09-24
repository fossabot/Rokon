package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/sethvargo/go-githubactions"
)

type Build struct {
	Chroots        []string `json:"chroots"`
	EndedOn        *int64   `json:"ended_on"`
	IsBackground   bool     `json:"is_background"`
	OwnerName      string   `json:"ownername"`
	ProjectDirName string   `json:"project_dirname"`
	ProjectName    string   `json:"projectname"`
	RepoURL        string   `json:"repo_url"`
	StartedOn      int64    `json:"started_on"`
	State          string   `json:"state"`
	SubmittedOn    int64    `json:"submitted_on"`
	Submitter      *string  `json:"submitter"`
	SourcePackage  struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		Version string `json:"version"`
	} `json:"source_package"`
	ID int64 `json:"id"`
}

type Response struct {
	Items  []Build `json:"items"`
	Output string
	Error  string
}

func main() {
	actions := githubactions.WithFieldsMap(map[string]string{
		"line": "100",
	})

	// Attempt to get COPR_TOKEN and COPR_LOGIN from GitHub Actions environment
	coprToken := actions.Getenv("COPR_TOKEN")
	if coprToken == "" {
		coprToken = os.Getenv("COPR_TOKEN") // Fallback to OS environment
		if coprToken == "" {
			log.Fatalf("missing 'COPR_TOKEN'")
		}
	}

	coprUsername := "brycensranch"

	// Attempt to get COPR_LOGIN from GitHub Actions environment
	coprLogin := actions.Getenv("COPR_LOGIN")
	if coprLogin == "" {
		coprLogin = os.Getenv("COPR_LOGIN") // Fallback to OS environment
		if coprLogin == "" {
			log.Fatalf("missing 'COPR_LOGIN'")
		}
	}

	projectName := "rokon"
	packageName := "rokon"

	// Fetch pending builds
	pendingBuilds := fetchBuilds(coprUsername, coprToken, coprLogin, projectName, packageName, "pending")

	// Fetch running builds
	runningBuilds := fetchBuilds(coprUsername, coprToken, coprLogin, projectName, packageName, "running")

	// Combine pending and running builds and exclude the latest one
	allBuilds := append(pendingBuilds, runningBuilds...)
	if len(allBuilds) > 0 {
		// Sort by SubmittedOn to find the latest build
		sort.Slice(allBuilds, func(i, j int) bool {
			return allBuilds[i].SubmittedOn < allBuilds[j].SubmittedOn
		})

		// Exclude the latest build from cancellation
		latestBuildID := allBuilds[len(allBuilds)-1].ID
		cancelBuilds(allBuilds, coprUsername, coprToken, coprLogin, latestBuildID)
	}
}

func fetchBuilds(coprUsername, coprToken, coprLogin, projectName, packageName, status string) []Build {
	// Create the URL for fetching builds based on the status
	apiURL := fmt.Sprintf("https://copr.fedorainfracloud.org/api_3/build/list?ownername=%s&projectname=%s&packagename=%s&status=%s", coprUsername, projectName, packageName, status)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Parse JSON response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var buildResponse Response
	if err := json.Unmarshal(body, &buildResponse); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return buildResponse.Items
}

func cancelBuilds(builds []Build, coprUsername, coprToken, coprLogin string, latestBuildID int64) {
	var wg sync.WaitGroup
	for _, build := range builds {
		if build.ID != latestBuildID && (build.State == "running" || build.State == "pending") {
			wg.Add(1) // Increment the WaitGroup counter
			go func(buildID int64) {
				defer wg.Done() // Decrement the counter when the goroutine completes
				cancelBuild(coprUsername, coprToken, coprLogin, buildID)
			}(build.ID)
		}
	}

	// Wait until all the goroutines for cancelling builds are done
	wg.Wait()
}

func cancelBuild(coprUsername string, coprToken string, coprLogin string, buildID int64) {
	cancelURL := fmt.Sprintf("https://copr.fedorainfracloud.org/api_3/build/cancel/%d", buildID)

	req, err := http.NewRequest(http.MethodPut, cancelURL, http.NoBody)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Encode login:token for basic auth (username:password)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", coprLogin, coprToken)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error cancelling build: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var cancelResponse Response
	if err := json.Unmarshal(body, &cancelResponse); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Build %d cancelled successfully.\n", buildID)
	} else {
		fmt.Printf("Failed to cancel build %d (%s). Status: %s\n", buildID, cancelURL, resp.Status)
		fmt.Printf("Headers: %v\n", req.Header)
		log.Fatalf("Message: %v\n", cancelResponse)
	}
}
