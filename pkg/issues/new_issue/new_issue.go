package new_issue

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/types"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/util"
)

type FilterFunc func(*github.Issue) bool

type GoodIssues struct {
	filepath string
	client   *github.Client
	// labels will be used as or condition
	labels []string

	shouldEvict FilterFunc
}

func NewDefaultFilter() FilterFunc {
	return func(issue *github.Issue) bool {
		return time.Since(*issue.CreatedAt) < 24*time.Hour
	}
}

func NewGoodIssues(
	filepath string,
	client *github.Client,
	labels []string,

	filterFunc FilterFunc,
) *GoodIssues {
	if filterFunc == nil {
		filterFunc = NewDefaultFilter()
	}

	return &GoodIssues{
		filepath:    filepath,
		client:      client,
		labels:      labels,
		shouldEvict: filterFunc,
	}
}

func (g *GoodIssues) GetAll(since time.Time) ([]*github.Issue, error) {
	f, err := os.Open(g.filepath)
	if err != nil {
		return nil, err
	}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(f)

	// Iterate over each line and store it in a slice of strings
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check for any errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	queries, err := util.ConvertToAPIEndpoint(lines)
	if err != nil {
		return nil, err
	}

	issues := []*github.Issue{}

	// loop through each labels because we want
	// or condition
	for _, l := range g.labels {
		issue := fromStringsToGithubIssues(
			context.TODO(),
			g.client,
			queries,
			&github.IssueListByRepoOptions{
				Labels: []string{l},
				Since:  since,
			},
		)
		if issue != nil {
			issues = append(issues, issue...)
		}
	}

	return g.filterIssues(issues), nil
}

func (g *GoodIssues) filterIssues(issues []*github.Issue) []*github.Issue {
	// Create a map to track the occurrence of each element
	occurrence := make(map[string]bool)
	uniqueSlice := []*github.Issue{}

	// Iterate over the slice
	for _, issue := range issues {
		if g.shouldEvict(issue) {
			continue
		}

		// Check if the element is already in the map
		if !occurrence[issue.GetHTMLURL()] {
			// Add the element to the map and the uniqueSlice
			occurrence[issue.GetHTMLURL()] = true
			uniqueSlice = append(uniqueSlice, issue)
		}
	}

	return uniqueSlice
}

func fromStringsToGithubIssues(
	ctx context.Context,
	client *github.Client,
	queries []types.IssueQueryWithNumber,
	opt *github.IssueListByRepoOptions,
) []*github.Issue {
	issues := []*github.Issue{}
	for _, q := range queries {
		issue, _, err := client.Issues.ListByRepo(ctx, q.Owner, q.Repo, opt)
		if err != nil {
			log.Println(err)
		} else {
			issues = append(issues, issue...)
		}
	}

	return issues
}
