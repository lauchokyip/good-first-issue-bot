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

type FilterFunc func(github.Issue) bool

type GoodIssues struct {
	filename string
	client   *github.Client
	labels   []string

	filterFunc FilterFunc
}

func NewDefaultFilter() FilterFunc {
	return func(issue github.Issue) bool {
		return time.Since(*issue.CreatedAt) < 24*time.Hour
	}
}

func NewGoodIssues(
	filename string,
	client *github.Client,
	labels []string,

	filterFunc FilterFunc,
) *GoodIssues {
	if filterFunc == nil {
		filterFunc = NewDefaultFilter()
	}

	return &GoodIssues{
		filename:   filename,
		client:     client,
		labels:     labels,
		filterFunc: filterFunc,
	}
}

func (g *GoodIssues) GetAll(since time.Time) ([]*github.Issue, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := dir + "/" + g.filename
	f, err := os.Open(path)
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

	issues := fromStringsToGithubIssues(
		context.TODO(),
		g.client,
		queries,
		&github.IssueListByRepoOptions{
			Labels: g.labels,
			Since:  since,
		},
	)

	return issues, nil
}

func fromStringsToGithubIssues(
	ctx context.Context,
	client *github.Client,
	queries []types.IssueQueryWithNumber,
	opt *github.IssueListByRepoOptions,
) []*github.Issue {
	issues := make([]*github.Issue, len(queries))
	for _, q := range queries {
		issue, _, err := client.Issues.ListByRepo(ctx, q.Owner, q.Repo, opt)
		if err != nil {
			log.Println(err)
		}
		issues = append(issues, issue...)
	}

	return issues
}
