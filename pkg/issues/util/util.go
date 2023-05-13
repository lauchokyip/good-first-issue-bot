package util

import (
	"strconv"
	"strings"

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
