# A Stupid Bot to scrape Github API for new good first issue and help wanted

* Use mailjet instead 
* This should get the issue
```go
// Issue represents a GitHub issue
type Issue struct {
	Title string `json:"title"`
	Url   string `json:"html_url"`
}

func main() {
	repoOwner := "<repository-owner>"
	repoName := "<repository-name>"

	// Make a request to the GitHub API to fetch the issues
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?labels=good%20first%20issue", repoOwner, repoName)
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Failed to fetch issues: %s\n", err.Error())
		return
	}
```

## Keep track of Subscribed Issues