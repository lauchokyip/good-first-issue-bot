package util

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/types"
)

func ConvertToAPIEndpoint(urls []string) ([]types.IssueQueryWithNumber, error) {
	queries := make([]types.IssueQueryWithNumber, len(urls))

	for i, url := range urls {
		// Remove the "https://github.com/" prefix
		url = strings.TrimPrefix(url, "https://github.com/")

		// Split the remaining URL into parts
		parts := strings.Split(url, "/")

		// Extract the owner, repo, and issue number
		owner := parts[0]
		repo := parts[1]

		var issueNumber int
		if len(parts) == 4 {
			var err error
			issueNumber, err = strconv.Atoi(parts[3])
			if err != nil {
				return nil, err
			}
		}

		queries[i] = types.IssueQueryWithNumber{
			IssueQuery: types.IssueQuery{
				Owner: owner,
				Repo:  repo,
			},
			Number: issueNumber,
		}
	}

	return queries, nil
}

func SlicesToMap(slices []string) map[string]bool {
	ret := map[string]bool{}

	for _, ss := range slices {
		ret[ss] = true
	}

	return ret
}

func GetLargestIssueNumber(ctx context.Context, client *github.Client, issue types.IssueQuery) (int, error) {
	issues, _, err := client.Issues.ListByRepo(ctx, issue.Owner, issue.Repo, &github.IssueListByRepoOptions{})
	if err != nil {
		return -1, err
	}

	// Find the largest issue number
	var largestIssueNumber int
	for _, issue := range issues {
		if issue.GetNumber() > largestIssueNumber {
			largestIssueNumber = issue.GetNumber()
		}
	}

	return largestIssueNumber, nil
}
